// server.go
package protosocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Server gerencia as conexões WebSocket.
type Server struct {
	upgrader     websocket.Upgrader
	clients      map[string]*Socket
	lock         sync.Mutex
	onConnection func(socket *Socket)
	handlers     map[string]func(proto.Message, *Socket)
}

// NewServer cria uma nova instância do Server.
func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			// Em produção, ajuste a checagem de origem conforme necessário.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:  make(map[string]*Socket),
		handlers: make(map[string]func(proto.Message, *Socket)),
	}
}

// OnConnection permite registrar um callback que será chamado quando um novo cliente se conectar.
func (s *Server) OnConnection(callback func(socket *Socket)) {
	s.onConnection = callback
}

// On permite registrar um handler para um evento específico.
func (s *Server) On(event string, handler func(proto.Message, *Socket)) {
	s.handlers[event] = handler
}

// Broadcast envia uma mensagem para todos os clientes conectados.
func (s *Server) Broadcast(event string, msg proto.Message) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, client := range s.clients {
		if err := client.Emit(event, msg); err != nil {
			log.Printf("Erro ao enviar mensagem para %s: %v\n", client.ID, err)
		}
	}
}

// ServeHTTP implementa o handler HTTP que fará o upgrade para WebSocket.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade da conexão:", err)
		return
	}

	socketID := uuid.New().String()
	socket := NewSocket(conn, socketID)

	s.lock.Lock()
	s.clients[socketID] = socket
	s.lock.Unlock()

	// Configura os handlers padrão
	for event, handler := range s.handlers {
		socket.On(event, handler)
	}

	if s.onConnection != nil {
		go s.onConnection(socket)
	}

	socket.Listen()

	s.lock.Lock()
	delete(s.clients, socketID)
	s.lock.Unlock()
}
