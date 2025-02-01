# ProtoSocket

ProtoSocket é uma biblioteca Go para comunicação WebSocket com suporte a Protocol Buffers, oferecendo recursos avançados de resiliência, monitoramento e segurança.

## Características

- 🔒 **Segurança**
  - Rate limiting
  - Firewall
  - Autenticação
  - Criptografia
  - Validação de mensagens

- 💪 **Resiliência**
  - Circuit breaker
  - Retry automático
  - Reconexão automática
  - Timeout handling
  - Buffering de mensagens

- 📊 **Monitoramento**
  - Health checks
  - Métricas
  - Tracing
  - Logging estruturado
  - Status em tempo real

- 🔌 **Extensibilidade**
  - Sistema de plugins
  - Handlers customizáveis
  - Configuração dinâmica
  - Middleware support

- 🚀 **Performance**
  - Protocol Buffers
  - Message sequencing
  - Priorização de mensagens
  - Connection pooling

## Instalação

```bash
go get github.com/mendes113/protosocket
```

## Uso Básico

```go
package main

import (
    "log"
    "github.com/mendes113/protosocket"
    "github.com/mendes113/protosocket/proto"
)

func main() {
    // Inicializa o cliente
    client := protosocket.NewClient("ws://localhost:8080")
    defer client.Close()

    // Configura handler de mensagens
    client.On("chat", func(msg proto.Message, c *protosocket.Client) {
        log.Printf("Mensagem recebida: %s", string(msg.Data))
    })

    // Envia mensagem
    msg := &proto.Message{
        Type: "chat",
        Data: []byte("Olá!"),
        Metadata: map[string]string{
            "sender": "exemplo",
        },
    }

    if err := client.Emit("chat", msg); err != nil {
        log.Printf("Erro: %v", err)
    }

    select {}
}
```

## Exemplos

O diretório `examples/` contém exemplos completos de uso:

- `sender/` - Cliente que envia mensagens periodicamente
- `receiver/` - Cliente que recebe e responde mensagens
- `advanced_client.go` - Cliente com recursos avançados

## Documentação

Para mais detalhes sobre cada componente:

- [Segurança](docs/security.md)
- [Resiliência](docs/resilience.md)
- [Monitoramento](docs/monitoring.md)
- [Plugins](docs/plugins.md)
- [Configuração](docs/config.md)

## Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/MinhaFeature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/MinhaFeature`)
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes. 