package protosocket

import (
	"log"
	"net/http"
)

type SocketServer interface {
	OnConnection(func(*Socket))
	Start(port string) error
}

type socketServer struct {
	*Server
	config ServerConfig
}

func New(config ...ServerConfig) SocketServer {
	var cfg ServerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return &socketServer{
		Server: NewServer(cfg),
		config: cfg,
	}
}

func (s *socketServer) Start(port string) error {
	http.Handle(s.config.Path, s.Server)
	log.Printf("Servidor WebSocket iniciado em %s%s", port, s.config.Path)
	return http.ListenAndServe(port, nil)
}
