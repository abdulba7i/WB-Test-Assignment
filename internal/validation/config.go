package validation

type ValidationConfig struct {
	MaxNameLength int
	MinAge        int
	MaxAge        int
}

func DefaultConfig() ValidationConfig {
	return ValidationConfig{
		MaxNameLength: 100,
		MinAge:        18,
		MaxAge:        100,
	}
}

var (
	globalConfig = DefaultConfig()
)

func SetGlobalConfig(cfg ValidationConfig) {
	globalConfig = cfg
}

func GetConfig() ValidationConfig {
	return globalConfig
}
