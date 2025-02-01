package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mendes113/protosocket/protosocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer origem
	},
}

var clients = make(map[string]*protosocket.Socket)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade:", err)
		return
	}

	clientID := uuid.New().String()
	socket := protosocket.NewSocket(conn, clientID)
	clients[clientID] = socket

	// Broadcast mensagens para todos os clientes
	socket.On("mensagem", func(data proto.Message, s *protosocket.Socket) {
		for _, client := range clients {
			if err := client.Emit("mensagem", data); err != nil {
				log.Printf("Erro ao enviar mensagem para %s: %v\n", client.ID, err)
			}
		}
	})

	go socket.Listen()
}

func main() {
	server := protosocket.NewServer()

	server.On("chat", func(data proto.Message, s *protosocket.Socket) {
		// Broadcast a mensagem para todos os clientes
		server.Broadcast("chat", data)
	})

	http.Handle("/ws", server)
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
