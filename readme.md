# Protosocket

Uma biblioteca Go para comunicação P2P usando WebSocket com suporte a Protocol Buffers.

## 🚀 Características

- ✨ Comunicação P2P bidirecional
- 🔒 Suporte a TLS/SSL
- 🚦 Rate limiting e controle de carga
- 📊 Telemetria e observabilidade
- 🔑 Middleware para autenticação
- 🗜️ Compressão de dados
- 📝 Logging estruturado (Zap)
- 🔄 Reconexão automática
- 🔍 Descoberta de serviços
- 📨 Garantia de ordem de mensagens

## 📦 Instalação

```bash
go get github.com/mendes113/protosocket
```

## 🎯 Uso Básico

### Iniciando um Peer

```go
package main

import (
    "github.com/mendes113/protosocket/protosocket"
    "google.golang.org/protobuf/proto"
)

func main() {
    // Cria um novo peer na porta 8081
    peer := protosocket.NewPeer(8081, "client1", "chat")

    // Conecta automaticamente a outros peers
    peer.AutoConnect([]int{8081, 8082})

    // Handler para mensagens de chat
    peer.On("chat", func(data proto.Message, senderID string) {
        msg := data.(*protosocket.ChatMessage)
        log.Printf("[%s]: %s\n", msg.Sender, msg.Content)
    })

    // Inicia o servidor
    peer.Start()
}
```

### Enviando Mensagens

```go
msg := &protosocket.ChatMessage{
    Content: "Olá, mundo!",
    Sender:  peer.ID,
}
peer.Broadcast("chat", msg)
```

### Transferência de Arquivos

```go
// Enviando arquivo
err := peer.SendBinary(targetID, "arquivo.txt", data, "text/plain")

// Recebendo arquivo
peer.ReceiveBinary(func(filename string, data []byte, sender string) {
    // Processa o arquivo recebido
})
```

## 🔒 Segurança

```go
config := protosocket.SecurityConfig{
    EnableTLS:      true,
    CertFile:       "cert.pem",
    KeyFile:        "key.pem",
    EnableCORS:     true,
    AllowedOrigins: []string{"localhost"},
}

peer.EnableSecurity(config)
```

## 🚦 Rate Limiting

```go
limiter := protosocket.NewRateLimiter(100, 10) // 100 req/s, burst 10
peer.Use(limiter.Middleware())
```

## 📊 Telemetria

```go
peer.EnableTelemetry(context.Background())
```

## 🔄 Reconexão Automática

```go
peer.EnableReconnect(protosocket.ReconnectConfig{
    MaxAttempts:       5,
    InitialDelay:      time.Second,
    MaxDelay:          time.Minute,
    BackoffMultiplier: 2.0,
})
```

## 🌐 Casos de Uso

- Chat distribuído
- Compartilhamento de arquivos P2P
- Jogos multiplayer descentralizados
- Sistemas de mensagens em tempo real
- Aplicações colaborativas

## 📋 Requisitos

- Go 1.22 ou superior
- Protocol Buffers
- WebSocket

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## ⚠️ Limitações Atuais

- Não suporta NAT traversal
- Sem suporte a DHT para descoberta
- Sem persistência de mensagens

## 📚 Documentação Adicional

Para mais detalhes sobre a API e exemplos, consulte a [documentação completa](docs/README.md).
```

Este README fornece uma visão geral completa da biblioteca, com exemplos de código e informações sobre recursos principais. A formatação com emojis torna a leitura mais agradável e ajuda na navegação visual do documento.

