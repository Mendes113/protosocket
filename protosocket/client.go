package protosocket

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	ID             string
	conn           *websocket.Conn
	handlers       map[string]func(proto.Message, *Client)
	logger         *zap.Logger
	sequence       uint64
	sequencer      *MessageSequencer
	metrics        *MetricsCollector
	circuitBreaker *CircuitBreaker
	validator      *MessageValidator
	retryConfig    RetryConfig
}

func NewClient(url string) *Client {
	logger := GetLogger()

	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		logger.Fatal("erro de conexão",
			zap.String("url", url),
			zap.Error(err))
	}

	c := &Client{
		ID:             uuid.New().String()[:8],
		conn:           conn,
		handlers:       make(map[string]func(proto.Message, *Client)),
		logger:         logger,
		sequencer:      NewMessageSequencer(),
		metrics:        NewMetricsCollector(),
		circuitBreaker: NewCircuitBreaker(5, time.Second*10),
		validator:      NewMessageValidator(),
		retryConfig: RetryConfig{
			MaxAttempts:       3,
			InitialDelay:      time.Second,
			MaxDelay:          time.Second * 5,
			BackoffMultiplier: 2.0,
		},
	}

	// Inicia a goroutine de escuta
	go c.listen()
	return c
}

func (c *Client) On(event string, handler func(proto.Message, *Client)) {
	c.handlers[event] = handler
}

func (c *Client) Emit(event string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	wrapper := &Message{
		Event: event,
		Data:  data,
	}

	wrapperData, err := proto.Marshal(wrapper)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.BinaryMessage, wrapperData)
}

func (c *Client) EmitWithRetry(event string, msg proto.Message) error {
	return WithRetry(c.retryConfig, func() error {
		return c.circuitBreaker.Execute(func() error {
			return c.Emit(event, msg)
		})
	})
}

func (c *Client) listen() {
	defer c.conn.Close()
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error("erro na leitura", zap.Error(err))
			return
		}

		c.metrics.RecordReceivedMessage(len(data))
		wrapper := &SequencedMessage{}
		if err := proto.Unmarshal(data, wrapper); err != nil {
			c.logger.Error("erro ao decodificar wrapper", zap.Error(err))
			continue
		}

		if err := c.processMessage(wrapper); err != nil {
			c.logger.Error("erro ao processar mensagem", zap.Error(err))
		}
	}
}

func (c *Client) processMessage(msg *SequencedMessage) error {
	// Validação
	if err := c.validator.Validate(msg); err != nil {
		c.logger.Error("mensagem inválida",
			zap.Error(err),
			zap.String("event", msg.Event),
			zap.String("sender", msg.SenderId))
		return err
	}

	// Processamento normal...
	return c.circuitBreaker.Execute(func() error {
		return c.processMessageInternal(msg)
	})
}

func (c *Client) processMessageInternal(msg *SequencedMessage) error {
	c.sequencer.lock.Lock()
	defer c.sequencer.lock.Unlock()

	lastSeq := c.sequencer.lastSeq[msg.SenderId]

	// Se é a próxima mensagem esperada
	if msg.Sequence == lastSeq+1 {
		c.sequencer.lastSeq[msg.SenderId] = msg.Sequence
		return c.deliverMessage(msg)
	}

	// Se é uma mensagem futura, guarda no buffer
	if msg.Sequence > lastSeq+1 {
		c.sequencer.buffer[msg.SenderId] = append(
			c.sequencer.buffer[msg.SenderId], *msg)
		c.tryDeliverBuffered(msg.SenderId)
	}

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Adicione os demais métodos (On, Emit, listen) aqui...
