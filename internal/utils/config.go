package utils

type Config struct {
	DatabaseUser     string `yaml:"databaseUser"`
	DatabasePassword string `yaml:"databasePassword"`
	DatabaseName     string `yaml:"databaseName"`
	DatabaseHost     string `yaml:"databaseHost"`
	DatabasePort     int    `yaml:"databasePort"`
	WebserverPort    int    `yaml:"webserverPort"`
}

func IsWebserverPortValid(cfg *Config) bool {
	if cfg.WebserverPort <= 0 || cfg.WebserverPort > 65535 {
		return false
	}
	return true
}
