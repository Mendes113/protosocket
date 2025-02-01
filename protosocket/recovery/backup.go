package recovery

import (
	"time"

	"github.com/mendes113/protosocket/protosocket/security"
)

type StorageProvider interface {
	Snapshot() ([]byte, error)
	Store(id string, data []byte) error
	Restore(id string) ([]byte, error)
	List() ([]string, error)
	Delete(id string) error
}

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type BackupManager struct {
	storage    StorageProvider
	schedule   *BackupSchedule
	retention  RetentionPolicy
	encryption security.Encryptor
}

type RecoveryPoint struct {
	ID        string
	Timestamp time.Time
	Size      int64
	Checksum  string
}
