package apiserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	DbURL    string `toml:"database_url"`
	SessionKey string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8000",
		LogLevel: "debug",
		DbURL: "host=localhost dbname=api_dev sslmode=disable",
	}
}
