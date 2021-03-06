package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config Represents database server and credentials
type Config struct {
	Server   string
	Database string
	Hostname string
	Port     string
}

// Read and parse the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
