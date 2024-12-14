package config

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ShortenerSettings struct {
	Length     int8 `yaml:"length"`
	MaxAttempt int8 `yaml:"max_attempt"`
}

type CacheConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database int    `yaml:"database"`
	TTL      int    `yaml:"ttl"`
	Username string `yaml:"username"` // TODO: not supported yet
	Password string `yaml:"password"` // TODO: not supported yet
}

type Config struct {
	Server         ServerConfig      `yaml:"server"`
	Db_Conn_String string            `yaml:"db_conn_string"`
	Cache          CacheConfig       `yaml:"cache"`
	ShortenerProps ShortenerSettings `yaml:"shortener_props"`
}
