package config

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type CacheConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database int    `yaml:"database"`
	TTL      int    `yaml:"ttl"`
	Username string `yaml:"username"` // TODO: not supported yet
	Password string `yaml:"password"` // TODO: not supported yet
}

type Config struct {
	Server         ServerConfig `yaml:"server"`
	Db_Conn_String string       `yaml:"db_conn_string"`
	Cache          CacheConfig  `yaml:"cache"`
}
