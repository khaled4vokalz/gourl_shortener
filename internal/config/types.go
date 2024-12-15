package config

type ServerConfig struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

type ShortenerSettings struct {
	Length     int8 `koanf:"length"`
	MaxAttempt int8 `koanf:"maxAttempt"`
}

type StorageConfig struct {
	Type           string `koanf:"type"`
	Db_Conn_String string `koanf:"dbConnString"`
}

type CacheConfig struct {
	Enabled  bool   `koanf:"enabled"`
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	Database int    `koanf:"database"`
	TTL      int    `koanf:"ttl"`
	Username string `koanf:"username"` // TODO: not supported yet
	Password string `koanf:"password"` // TODO: not supported yet
}

type Config struct {
	Server         ServerConfig      `koanf:"server"`
	Cache          CacheConfig       `koanf:"cache"`
	Storage        StorageConfig     `koanf:"storage"`
	ShortenerProps ShortenerSettings `koanf:"shortenerProps"`
}
