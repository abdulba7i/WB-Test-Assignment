package config

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  struct {
		Address     string `yaml:"address"`
		Timeout     string `yaml:"timeout"`
		IdleTimeout string `yaml:"idle_timeout"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
}

func NewConfig() *Config {
	return &Config{}
}
