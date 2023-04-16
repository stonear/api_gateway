package db

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Host     string `envconfig:"DB_HOST" default:""`
	Port     int    `envconfig:"DB_PORT" default:""`
	Database string `envconfig:"DB_DATABASE" default:""`
	Username string `envconfig:"DB_USERNAME" default:""`
	Password string `envconfig:"DB_PASSWORD" default:""`
}

func getConfig() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
