package protosocket

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type SecurityConfig struct {
	TLSConfig      *tls.Config
	EnableTLS      bool
	CertFile       string
	KeyFile        string
	AllowedOrigins []string
	EnableCORS     bool
}

func (p *Peer) EnableSecurity(config SecurityConfig) error {
	// Configuração de TLS
	if config.EnableTLS {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return err
		}

		// Modifique o método Start do Peer para usar TLS
		p.startServer = func() error {
			server := &http.Server{
				Addr: fmt.Sprintf(":%d", p.Port),
				TLSConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
					MinVersion:   tls.VersionTLS12,
				},
			}
			http.HandleFunc("/ws", p.handleWebSocket)
			return server.ListenAndServeTLS("", "")
		}
	}

	// Configuração de CORS
	if config.EnableCORS {
		p.upgrader.CheckOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowed := range config.AllowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		}
	}

	return nil
}
