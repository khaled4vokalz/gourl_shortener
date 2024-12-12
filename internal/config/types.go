package config

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Domain string       `yaml:"domain"`
}
