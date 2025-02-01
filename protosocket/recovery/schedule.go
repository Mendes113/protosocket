package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type BackupSchedule struct {
	Interval   time.Duration
	RetryCount int
	RetryDelay time.Duration
	LastBackup time.Time
	NextBackup time.Time
}

type RetentionPolicy struct {
	MaxAge     time.Duration
	MaxSize    int64
	MaxBackups int
	KeepLatest int
}

func (bm *BackupManager) CreateBackup() (*RecoveryPoint, error) {
	data, err := bm.storage.Snapshot()
	if err != nil {
		return nil, err
	}

	encrypted, err := bm.encryption.Encrypt(data)
	if err != nil {
		return nil, err
	}

	point := &RecoveryPoint{
		ID:        generateID(),
		Timestamp: time.Now(),
		Size:      int64(len(encrypted)),
		Checksum:  calculateChecksum(encrypted),
	}

	if err := bm.storage.Store(point.ID, encrypted); err != nil {
		return nil, err
	}

	return point, nil
}

func generateID() string {
	return uuid.New().String()[:8]
}

func calculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
