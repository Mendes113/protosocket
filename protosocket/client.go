package protosocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn     *websocket.Conn
	handlers map[string]func(proto.Message, *Client)
	Listen   func()
}

func NewClient(url string) *Client {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Fatal("Erro de conexão:", err)
	}

	c := &Client{
		conn:     conn,
		handlers: make(map[string]func(proto.Message, *Client)),
	}

	go c.Listen()
	return c
}

// Adicione os demais métodos (On, Emit, listen) aqui...
