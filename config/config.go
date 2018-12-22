package config

// Config Represents database server and credentials
type Config struct {
	Server   string
	Database string
	Hostname string
	Port     string
}

// // Read and parse the configuration file
// func (c *Config) Read() {
// 	c.Server = os.Getenv("server")
// 	c.Database = os.Getenv("database")
// 	c.Hostname = os.Getenv("hostname")

// 	if c.Database != "" {
// 		return
// 	}

// 	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
// 		log.Fatal(err)
// 	}
// }
