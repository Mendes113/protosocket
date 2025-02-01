package config

type ConfigManager struct {
	current   *Config
	history   []*ConfigChange
	watchers  []ConfigWatcher
	validator ConfigValidator
}

type ConfigValidator interface {
	Validate(updates map[string]interface{}) error
}

type ConfigWatcher interface {
	OnChange(change *ConfigChange)
}
