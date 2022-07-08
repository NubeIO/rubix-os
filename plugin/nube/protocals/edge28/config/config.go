package config

type Config struct {
	EnablePolling bool   `yaml:"enable_polling"`
	LogLevel      string `yaml:"log_level"`
	PollRate      int    `yaml:"poll_rate"`
}
