package protosocket

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pb "github.com/mendes113/protosocket/protosocket/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type websocketMessage struct {
	data []byte
	err  error
}

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
	textHandlers   map[string]func(string, *Client)
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
		textHandlers: make(map[string]func(string, *Client)),
	}

	// Inicia a goroutine de escuta
	go c.listen()
	return c
}

func (c *Client) On(event string, handler func(proto.Message, *Client)) {
	c.handlers[event] = handler
}

func (c *Client) Emit(event string, msg proto.Message) error {
	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	wrapper := &pb.MessageWrapper{
		Event:    event,
		Data:     payload,
		SenderId: c.ID,
		Sequence: atomic.AddUint64(&c.sequence, 1),
	}

	data, err := proto.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("erro ao serializar wrapper: %w", err)
	}

	return c.conn.WriteMessage(websocket.BinaryMessage, data)
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
		msgType, data, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error("erro na leitura", zap.Error(err))
			return
		}

		c.metrics.RecordReceivedMessage(len(data))

		switch msgType {
		case websocket.BinaryMessage:
			// Processa mensagens Protobuf
			var wrapper pb.MessageWrapper
			if err := proto.Unmarshal(data, &wrapper); err != nil {
				c.logger.Error("erro ao decodificar wrapper", zap.Error(err))
				continue
			}

			if handler, ok := c.handlers[wrapper.Event]; ok {
				var payload pb.Message
				if err := proto.Unmarshal(wrapper.Data, &payload); err != nil {
					c.logger.Error("erro ao decodificar payload", zap.Error(err))
					continue
				}
				handler(&payload, c)
			}

		case websocket.TextMessage:
			// Processa mensagens de texto simples
			if textHandler, ok := c.textHandlers["text"]; ok {
				textHandler(string(data), c)
			}

		default:
			c.logger.Warn("tipo de mensagem não suportado", zap.Int("tipo", msgType))
		}
	}
}

func (c *Client) processMessage(msg *websocketMessage) error {
	var wrapper pb.MessageWrapper
	if err := proto.Unmarshal(msg.data, &wrapper); err != nil {
		return fmt.Errorf("erro ao decodificar wrapper: %v", err)
	}

	// Verifica se a mensagem é para este cliente
	if wrapper.SenderId != "" && wrapper.SenderId != c.ID {
		return nil // Ignora mensagens de outros clientes
	}

	if handler, ok := c.handlers[wrapper.Event]; ok {
		var payload pb.Message
		if err := proto.Unmarshal(wrapper.Data, &payload); err != nil {
			return fmt.Errorf("erro ao decodificar payload: %v", err)
		}
		handler(&payload, c)
		return nil
	}

	return fmt.Errorf("handler não encontrado para o evento: %s", wrapper.Event)
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

// Novo método para registrar handlers de texto
func (c *Client) OnText(handler func(string, *Client)) {
	c.textHandlers["text"] = handler
}

// Novo método para enviar mensagens de texto
func (c *Client) EmitText(text string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(text))
}

// Adicione os demais métodos (On, Emit, listen) aqui...
