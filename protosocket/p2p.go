package protosocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Peer struct {
	ID          string
	Port        int
	ServiceInfo *ServiceInfo
	upgrader    websocket.Upgrader
	clients     map[string]*Client
	lock        sync.RWMutex
	handlers    map[string]func(proto.Message, string)
	discovery   *ServiceDiscovery
	telemetry   *Telemetry
	startServer func() error
	logger      *zap.Logger
}

type ServiceDiscovery struct {
	services map[string]*ServiceInfo
	lock     sync.RWMutex
}

func NewPeer(port int, serviceName, serviceType string) *Peer {
	peer := &Peer{
		ID:   uuid.New().String()[:8],
		Port: port,
		ServiceInfo: &ServiceInfo{
			Id:   uuid.New().String()[:8],
			Name: serviceName,
			Type: serviceType,
			Metadata: map[string]string{
				"host": fmt.Sprintf("localhost:%d", port),
			},
		},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:  make(map[string]*Client),
		handlers: make(map[string]func(proto.Message, string)),
		discovery: &ServiceDiscovery{
			services: make(map[string]*ServiceInfo),
		},
		logger: GetLogger(),
	}

	// Registra handlers de descoberta
	peer.On("service.discover", peer.handleServiceDiscover)
	peer.On("service.announce", peer.handleServiceAnnounce)

	return peer
}

func (p *Peer) handleServiceDiscover(msg proto.Message, senderID string) {
	// Envia informações do serviço local
	p.Broadcast("service.announce", &ChatMessage{
		Type:    MessageType_SERVICE,
		Service: p.ServiceInfo,
	})
}

func (p *Peer) handleServiceAnnounce(msg proto.Message, senderID string) {
	if cm, ok := msg.(*ChatMessage); ok && cm.Service != nil {
		p.discovery.lock.Lock()
		p.discovery.services[cm.Service.Id] = cm.Service
		p.discovery.lock.Unlock()

		log.Printf("Novo serviço descoberto: %s (%s)\n", cm.Service.Name, cm.Service.Type)
	}
}

// Encontra serviços por tipo
func (p *Peer) FindServices(serviceType string) []*ServiceInfo {
	p.discovery.lock.RLock()
	defer p.discovery.lock.RUnlock()

	var services []*ServiceInfo
	for _, service := range p.discovery.services {
		if service.Type == serviceType {
			services = append(services, service)
		}
	}
	return services
}

// Conecta a um serviço específico
func (p *Peer) ConnectToService(service *ServiceInfo) error {
	if host, ok := service.Metadata["host"]; ok {
		return p.Connect(host)
	}
	return fmt.Errorf("host não encontrado para o serviço %s", service.Id)
}

// Conecta a outro peer
func (p *Peer) Connect(addr string) error {
	client := NewClient(fmt.Sprintf("ws://%s/ws", addr))

	p.logger.Info("conectando ao peer",
		zap.String("addr", addr),
		zap.String("clientID", client.ID))

	// Handler genérico para repassar mensagens
	messageHandler := func(msg proto.Message, c *Client) {
		if handler, ok := p.handlers["chat"]; ok {
			handler(msg, c.ID)
		}
		p.Broadcast("chat", msg)
	}

	client.On("chat", messageHandler)
	client.On("binary", func(msg proto.Message, c *Client) {
		if handler, ok := p.handlers["binary"]; ok {
			handler(msg, c.ID)
		}
	})

	p.lock.Lock()
	p.clients[client.ID] = client
	p.lock.Unlock()

	go client.listen()

	p.logger.Info("conexão estabelecida",
		zap.String("addr", addr),
		zap.String("clientID", client.ID))

	return nil
}

// Registra um handler para eventos
func (p *Peer) On(event string, handler func(proto.Message, string)) {
	p.handlers[event] = handler
}

// Envia mensagem para todos os peers conectados
func (p *Peer) Broadcast(event string, msg proto.Message) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	senderID := ""
	if cm, ok := msg.(*ChatMessage); ok {
		senderID = cm.Sender
		event = "chat" // Força evento "chat" para mensagens de texto
	}

	for _, client := range p.clients {
		if client.ID != senderID {
			if err := client.Emit(event, msg); err != nil {
				log.Printf("Erro ao enviar para peer %s: %v\n", client.ID, err)
			}
		}
	}
}

func (p *Peer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade:", err)
		return
	}

	client := &Client{
		ID:       uuid.New().String()[:8],
		conn:     conn,
		handlers: make(map[string]func(proto.Message, *Client)),
	}

	// Configura os handlers para o novo cliente
	for event, handler := range p.handlers {
		finalHandler := handler // Captura o handler atual
		client.On(event, func(msg proto.Message, c *Client) {
			finalHandler(msg, c.ID)
			// Repassa a mensagem para outros clientes
			p.Broadcast(event, msg)
		})
	}

	p.lock.Lock()
	p.clients[client.ID] = client
	p.lock.Unlock()

	go client.listen()
}

// Inicia o peer
func (p *Peer) Start() error {
	if p.startServer != nil {
		return p.startServer()
	}
	http.HandleFunc("/ws", p.handleWebSocket)
	return http.ListenAndServe(fmt.Sprintf(":%d", p.Port), nil)
}

// Modifique o método SendBinary para incluir um tipo específico
func (p *Peer) SendBinary(targetID string, filename string, data []byte, mimeType string) error {
	log.Printf("Enviando arquivo para %s. Tamanho: %d bytes", targetID, len(data))
	msg := &BinaryMessage{
		Filename:  filename,
		Content:   data,
		Size:      int64(len(data)),
		MimeType:  mimeType,
		Sender:    p.ID,
		Timestamp: time.Now().Unix(),
		Type:      MessageType_BINARY,
	}

	p.lock.RLock()
	client, exists := p.clients[targetID]
	p.lock.RUnlock()

	if !exists {
		return fmt.Errorf("peer %s não encontrado", targetID)
	}

	return client.Emit("binary", msg) // Usa evento específico "binary"
}

// Adicione este método ao Peer para registrar handlers de binário corretamente
func (p *Peer) OnBinary(handler func(filename string, data []byte, sender string)) {
	p.On("binary", func(msg proto.Message, senderID string) {
		if binaryMsg, ok := msg.(*BinaryMessage); ok {
			handler(binaryMsg.Filename, binaryMsg.Content, binaryMsg.Sender)
		}
	})
}

// Modifique o método ReceiveBinary para usar o novo handler
func (p *Peer) ReceiveBinary(handler func(filename string, data []byte, sender string)) {
	p.OnBinary(handler)
}

// GetClients retorna uma lista de IDs dos clients conectados
func (p *Peer) GetClients() []string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var clientIDs []string
	for id := range p.clients {
		clientIDs = append(clientIDs, id)
	}
	return clientIDs
}

func (p *Peer) IsConnected() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return len(p.clients) > 0
}

func (p *Peer) AutoConnect(otherPorts []int) {
	go func() {
		for {
			if !p.IsConnected() {
				for _, port := range otherPorts {
					if port == p.Port {
						continue
					}
					addr := fmt.Sprintf("localhost:%d", port)
					if err := p.Connect(addr); err == nil {
						log.Printf("Conectado com sucesso ao peer na porta %d", port)
						break
					}
				}
			}
			time.Sleep(time.Second)
		}
	}()
}
