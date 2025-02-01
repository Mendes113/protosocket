// main.go
package main

import (
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/mendes113/protosocket/protosocket"
)

type SocketHandler interface {
	On(event string, handler func(data proto.Message, socket *protosocket.Socket))
	Emit(event string, data proto.Message) error
}

type ServerConfig struct {
	Port        string
	Path        string
	PingHandler func() proto.Message
}

func main() {
	server := protosocket.NewServer()

	server.OnConnection(func(socket *protosocket.Socket) {
		log.Println("Novo cliente conectado:", socket.ID)

		// Handler tipado
		protosocket.OnTyped[*protosocket.SimpleMessage](socket, "ping",
			protosocket.EventHandler[*protosocket.SimpleMessage]{
				Handler: func(msg *protosocket.SimpleMessage, s *protosocket.Socket) {
					log.Printf("Recebido: %s", msg.Text)

					resp := &protosocket.SimpleMessage{Text: "pong!"}
					if err := s.Emit("pong", resp); err != nil {
						log.Println("Erro:", err)
					}
				},
			},
		)
	})

	log.Fatal(http.ListenAndServe(":8000", server))
}
