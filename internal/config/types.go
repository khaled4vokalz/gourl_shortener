package config

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	Server         ServerConfig `yaml:"server"`
	Db_Conn_String string       `yaml:"db_conn_string"`
}
