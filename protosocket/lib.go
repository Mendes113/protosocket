package protosocket

import (
	"log"
	"net/http"
)

type ServerConfig struct {
	Path string
}



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
		Server: NewServer(),
		config: cfg,
	}
}

func (s *socketServer) Start(port string) error {
	if s.config.Path == "" {
		s.config.Path = "/"
	} else if s.config.Path[0] != '/' {
		s.config.Path = "/" + s.config.Path
	}

	http.Handle(s.config.Path, s.Server)
	log.Printf("Servidor WebSocket iniciado em %s%s", port, s.config.Path)
	return http.ListenAndServe(port, nil)
}
