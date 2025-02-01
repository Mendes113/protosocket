// socket.go
package protosocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Socket representa uma conexão com um cliente.
type Socket struct {
	Conn         *websocket.Conn
	ID           string
	events       map[string]func(data proto.Message, socket *Socket)
	lock         sync.Mutex
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewSocket cria um novo Socket com o ID fornecido.
func NewSocket(conn *websocket.Conn, id string) *Socket {
	return &Socket{
		Conn:   conn,
		ID:     id,
		events: make(map[string]func(data proto.Message, socket *Socket)),
	}
}

// On registra um handler para um evento específico.
func (s *Socket) On(event string, callback func(data proto.Message, socket *Socket)) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.events[event] = callback
}

// Emit envia uma mensagem para o cliente usando protobuf.
// O parâmetro `data` deve ser uma mensagem protobuf que será empacotada no campo Any.
func (s *Socket) Emit(event string, data proto.Message) error {
	anyData, err := anypb.New(data)
	if err != nil {
		return err
	}

	// Cria a mensagem protobuf.
	msg := &Message{
		Event: event,
		Data:  anyData,
	}

	// Serializa a mensagem para []byte.
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Envia a mensagem como binário.
	return s.Conn.WriteMessage(websocket.BinaryMessage, b)
}

// SetTimeouts define os tempos limite para leitura e escrita.
func (s *Socket) SetTimeouts(read, write time.Duration) {
	s.readTimeout = read
	s.writeTimeout = write
}

// Listen fica em loop lendo mensagens do cliente e invoca os handlers registrados.
func (s *Socket) Listen() {
	defer s.Conn.Close()
	for {
		if s.readTimeout > 0 {
			s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		}

		msgType, b, err := s.Conn.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			break
		}

		// Verifica se a mensagem é binária.
		if msgType != websocket.BinaryMessage {
			log.Println("Mensagem recebida não é binária; ignorando")
			continue
		}

		// Desserializa a mensagem usando protobuf.
		var msg Message
		if err := proto.Unmarshal(b, &msg); err != nil {
			log.Println("Erro ao desserializar mensagem:", err)
			continue
		}

		// Extrai o dado do Any
		var anyData anypb.Any
		if err := msg.Data.UnmarshalTo(&anyData); err != nil {
			log.Println("Erro ao decodificar Any:", err)
			continue
		}

		// Desserializa para o tipo concreto
		concreteMsg, err := anyData.UnmarshalNew()
		if err != nil {
			log.Println("Erro ao desserializar mensagem concreta:", err)
			continue
		}

		s.lock.Lock()
		handler, exists := s.events[msg.Event]
		s.lock.Unlock()

		if exists && handler != nil {
			go handler(concreteMsg, s)
		} else {
			log.Printf("Nenhum handler registrado para o evento '%s'\n", msg.Event)
		}
	}
}
