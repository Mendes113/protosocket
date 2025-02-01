package protosocket

import (
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidMessage  = errors.New("mensagem inválida")
	ErrMessageTooLarge = errors.New("mensagem muito grande")
	ErrInvalidSender   = errors.New("remetente inválido")
	ErrMessageExpired  = errors.New("mensagem expirada")
)

type Validator interface {
	Validate(msg proto.Message) error
}

type MessageValidator struct {
	maxMessageSize    int
	messageTimeout    time.Duration
	allowedEventTypes map[string]bool
}

func NewMessageValidator() *MessageValidator {
	return &MessageValidator{
		maxMessageSize: 1024 * 1024, // 1MB
		messageTimeout: time.Minute * 5,
		allowedEventTypes: map[string]bool{
			"chat":    true,
			"binary":  true,
			"service": true,
		},
	}
}

func (v *MessageValidator) Validate(msg *SequencedMessage) error {
	// Validação básica
	if msg == nil {
		return ErrInvalidMessage
	}

	// Tamanho da mensagem
	if len(msg.Data) > v.maxMessageSize {
		return ErrMessageTooLarge
	}

	// Tipo de evento
	if !v.allowedEventTypes[msg.Event] {
		return ErrInvalidMessage
	}

	// Timestamp
	msgTime := time.Unix(msg.Timestamp, 0)
	if time.Since(msgTime) > v.messageTimeout {
		return ErrMessageExpired
	}

	return nil
}
