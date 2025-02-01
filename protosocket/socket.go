// socket.go
package protosocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
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
	// Serializa o dado diretamente para bytes
	msgData, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// Cria a mensagem protobuf
	msg := &Message{
		Event: event,
		Data:  msgData, // Agora usa []byte diretamente
	}

	// Serializa a mensagem para []byte
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

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

		if msgType != websocket.BinaryMessage {
			log.Println("Mensagem recebida não é binária; ignorando")
			continue
		}

		var wrapper Message
		if err := proto.Unmarshal(b, &wrapper); err != nil {
			log.Println("Erro ao desserializar wrapper:", err)
			continue
		}

		// Desserializa para o tipo correto baseado no evento
		var msg proto.Message
		switch wrapper.Event {
		case "chat":
			msg = &ChatMessage{}
		case "binary":
			msg = &BinaryMessage{}
		default:
			log.Printf("Evento desconhecido: %s\n", wrapper.Event)
			continue
		}

		if err := proto.Unmarshal(wrapper.Data, msg); err != nil {
			log.Println("Erro ao desserializar mensagem concreta:", err)
			continue
		}

		s.lock.Lock()
		handler, exists := s.events[wrapper.Event]
		s.lock.Unlock()

		if exists && handler != nil {
			go handler(msg, s)
		} else {
			log.Printf("Nenhum handler registrado para '%s'\n", wrapper.Event)
		}
	}
}
