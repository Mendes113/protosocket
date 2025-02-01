// server.go
package protosocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Adicione esta struct acima da definição do Server
type ServerConfig struct {
	Port         string
	Path         string
	CheckOrigin  func(r *http.Request) bool
	OnConnection func(socket *Socket)
	OnMessage    func(socket *Socket, message []byte)
	OnClose      func(socket *Socket)
	OnError      func(socket *Socket, err error)
	OnPing       func(socket *Socket)
	OnPong       func(socket *Socket)
	PingHandler  func() proto.Message
	Middlewares  []func(*Socket) error
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Server gerencia as conexões WebSocket.
type Server struct {
	upgrader     websocket.Upgrader
	clients      map[string]*Socket
	lock         sync.Mutex
	onConnection func(socket *Socket)
	config       ServerConfig
}

// NewServer cria uma nova instância do Server.
func NewServer(cfg ...ServerConfig) *Server {
	config := ServerConfig{
		Port:        ":8080",
		Path:        "/ws",
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	if len(cfg) > 0 {
		config = cfg[0]
	}

	return &Server{
		upgrader: websocket.Upgrader{
			// Em produção, ajuste a checagem de origem conforme necessário.
			CheckOrigin: config.CheckOrigin,
		},
		clients: make(map[string]*Socket),
		config:  config,
	}
}

// OnConnection permite registrar um callback que será chamado quando um novo cliente se conectar.
func (s *Server) OnConnection(callback func(socket *Socket)) {
	s.onConnection = callback
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
	socket.SetTimeouts(s.config.ReadTimeout, s.config.WriteTimeout)

	// Executa middlewares
	for _, middleware := range s.config.Middlewares {
		if err := middleware(socket); err != nil {
			conn.Close()
			return
		}
	}

	s.lock.Lock()
	s.clients[socketID] = socket
	s.lock.Unlock()

	if s.onConnection != nil {
		go s.onConnection(socket)
	}

	socket.Listen()

	s.lock.Lock()
	delete(s.clients, socketID)
	s.lock.Unlock()
}
