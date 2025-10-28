package main

import (
	"fmt"
	"net/url"
)

type config struct {
	Database   Database `envPrefix:"DB_"`
	Cache      Cache    `envPrefix:"CACHE_"`
	WebhookURL string   `env:"WEBHOOK_URL"`
	ServerPort int      `env:"SERVER_PORT"`
}
type Database struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Database string `env:"NAME"`
	SSLMode  string `env:"SSL_MODE"`
}

func (d Database) ConnectionString() string {
	query := url.Values{}
	query.Add("sslmode", d.SSLMode)

	ds := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.Username, d.Password),
		Host:     fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:     d.Database,
		RawQuery: query.Encode(),
	}
	return ds.String()
}

type Cache struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}
