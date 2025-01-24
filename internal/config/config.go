package config

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL: "host=localhost dbname=glouton_db sslmode=disable", // ces 3 arguments permettent de configurer le serveur
		ServerPort:  ":8080",
	}
}
