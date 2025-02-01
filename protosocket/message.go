package protosocket

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MessageSequencer struct {
	lastSeq     map[string]uint64 // Por remetente
	buffer      map[string][]SequencedMessage
	lock        sync.RWMutex
	maxBuffer   int
	maxWaitTime time.Duration
}

func NewMessageSequencer() *MessageSequencer {
	return &MessageSequencer{
		lastSeq:     make(map[string]uint64),
		buffer:      make(map[string][]SequencedMessage),
		maxBuffer:   1000,
		maxWaitTime: 5 * time.Second,
	}
}

// EmitSequenced envia uma mensagem com garantia de ordem
func (c *Client) EmitSequenced(event string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	wrapper := &SequencedMessage{
		Event:     event,
		Data:      data,
		Sequence:  atomic.AddUint64(&c.sequence, 1),
		Timestamp: time.Now().Unix(),
		SenderId:  c.ID,
	}

	wrapperData, err := proto.Marshal(wrapper)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.BinaryMessage, wrapperData)
}

func (c *Client) tryDeliverBuffered(senderID string) {
	messages := c.sequencer.buffer[senderID]
	lastSeq := c.sequencer.lastSeq[senderID]

	// Ordena mensagens por sequência
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Sequence < messages[j].Sequence
	})

	// Entrega mensagens em ordem
	for i, msg := range messages {
		if msg.Sequence == lastSeq+1 {
			c.deliverMessage(&messages[i])
			c.sequencer.lastSeq[senderID] = msg.Sequence
			lastSeq = msg.Sequence
		} else {
			break
		}
	}

	// Remove mensagens entregues do buffer
	c.sequencer.buffer[senderID] = messages[len(messages):]
}

func (c *Client) deliverMessage(msg *SequencedMessage) error {
	var protoMsg proto.Message
	switch msg.Event {
	case "chat":
		protoMsg = &ChatMessage{}
	case "binary":
		protoMsg = &BinaryMessage{}
	default:
		return fmt.Errorf("evento não suportado: %s", msg.Event)
	}

	if err := proto.Unmarshal(msg.Data, protoMsg); err != nil {
		return err
	}

	if handler, exists := c.handlers[msg.Event]; exists {
		handler(protoMsg, c)
	}

	return nil
}
