package utils

type Config struct {
	DatabaseUser     string `yaml:"databaseUser"`
	DatabasePassword string `yaml:"databasePassword"`
	DatabaseName     string `yaml:"databaseName"`
	DatabaseHost     string `yaml:"databaseHost"`
	DatabasePort     int    `yaml:"databasePort"`
	WebserverPort    int    `yaml:"webserverPort"`
}
