package protosocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	ID       string
	conn     *websocket.Conn
	handlers map[string]func(proto.Message, *Client)
	logger   *zap.Logger
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
		ID:       uuid.New().String()[:8],
		conn:     conn,
		handlers: make(map[string]func(proto.Message, *Client)),
		logger:   logger,
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

func (c *Client) listen() {
	defer c.conn.Close()
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Erro na leitura:", err)
			return
		}

		wrapper := &Message{}
		if err := proto.Unmarshal(data, wrapper); err != nil {
			log.Println("Erro ao decodificar wrapper:", err)
			continue
		}

		var msg proto.Message
		switch wrapper.Event {
		case "chat":
			msg = &ChatMessage{}
		case "binary":
			msg = &BinaryMessage{}
		default:
			log.Printf("Evento não suportado: %s\n", wrapper.Event)
			continue
		}

		if err := proto.Unmarshal(wrapper.Data, msg); err != nil {
			log.Println("Erro ao decodificar payload:", err)
			continue
		}

		if handler, exists := c.handlers[wrapper.Event]; exists {
			handler(msg, c)
		}
	}
}

// Adicione os demais métodos (On, Emit, listen) aqui...
