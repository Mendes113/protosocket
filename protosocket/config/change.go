package config

import (
	"time"
)

type ConfigChange struct {
	ID        string
	Timestamp time.Time
	OldValue  interface{}
	NewValue  interface{}
	User      string
}

type Config struct {
	Version     string                 `json:"version"`
	Settings    map[string]interface{} `json:"settings"`
	LastUpdated time.Time              `json:"lastUpdated"`
}

func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	if err := cm.validator.Validate(updates); err != nil {
		return err
	}

	change := &ConfigChange{
		Timestamp: time.Now(),
		OldValue:  cm.current.Settings,
		NewValue:  updates,
	}

	cm.current.Settings = updates
	cm.current.LastUpdated = time.Now()
	cm.history = append(cm.history, change)

	for _, watcher := range cm.watchers {
		go watcher.OnChange(change)
	}

	return nil
}
