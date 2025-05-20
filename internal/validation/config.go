package validation

// ValidationConfig содержит настройки валидации
type ValidationConfig struct {
	MaxNameLength int
	MinAge        int
	MaxAge        int
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() ValidationConfig {
	return ValidationConfig{
		MaxNameLength: 100,
		MinAge:        18,
		MaxAge:        100,
	}
}

var (
	// глобальная конфигурация
	globalConfig = DefaultConfig()
)

// SetGlobalConfig устанавливает глобальную конфигурацию
func SetGlobalConfig(cfg ValidationConfig) {
	globalConfig = cfg
}

// GetConfig возвращает текущую конфигурацию
func GetConfig() ValidationConfig {
	return globalConfig
}
