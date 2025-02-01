package security

import (
	"net"
)

type AuthProvider interface {
	Authenticate(token string) (bool, error)
	ValidateCredentials(username, password string) (string, error)
	RevokeToken(token string) error
}

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type AuditLogger interface {
	LogAccess(event string, success bool)
	LogSecurity(event string, severity string)
}

type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
	GetLimit() int
	GetRemaining(key string) int
}

type SecurityManager struct {
	authProvider AuthProvider
	encryptor    Encryptor
	rateLimiter  RateLimiter
	firewall     Firewall
	auditor      AuditLogger
}

type Firewall struct {
	blacklist map[string]bool
	rules     []FirewallRule
	ipRanges  []*net.IPNet
}

func (sm *SecurityManager) CheckRateLimit(key string) bool {
	if sm.rateLimiter == nil {
		return true
	}
	return sm.rateLimiter.Allow(key)
}

func (sm *SecurityManager) ResetRateLimit(key string) {
	if sm.rateLimiter != nil {
		sm.rateLimiter.Reset(key)
	}
}
