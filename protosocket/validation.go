package protosocket

import "google.golang.org/protobuf/proto"

type Validator interface {
	Validate(msg proto.Message) error
}

type MessageValidator struct {
	maxSize      int64
	allowedTypes map[string]bool
	sanitizers   []func([]byte) []byte
}

func (mv *MessageValidator) Validate(msg proto.Message) error {
	// Implementa validações de segurança
	return nil
}
