package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/utils"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const ENV_PREFIX = "GOURLAPP_"

var singleton *Config

func getPath() string {
	common.LoadEnv()
	config_path := os.Getenv("CONFIG_PATH")

	if config_path == "" {
		config_path = "configuration/dev.yaml"
	}
	return config_path
}

// TODO: this approach is not good, now anyone will be able to change this.
// this was done like this, so that it can be mocked in the tests
var GetConfig = func() *Config {
	if singleton == nil {
		LoadConfig()
	}
	return singleton
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")
	var config Config
	filePath := getPath()
	if err := k.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("Error loading YAML file: %v\n", err)
	}

	k.Load(env.Provider(ENV_PREFIX, ".", func(s string) string {
		return strings.Replace(
			strings.TrimPrefix(s, ENV_PREFIX), "_", ".", -1)
	}), nil)

	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("Error unmarshaling config: %v\n", err)
	}
	if config.Environment == "" {
		config.Environment = "prod"
	}
	if config.UrlsExpiresIn != "" {
		config.UrlsExpiresAt, _ = utils.ParseTime(config.UrlsExpiresIn)
	} else {
		config.UrlsExpiresAt, _ = utils.ParseTime("30d")
	}
	singleton = &config

	return singleton, nil
}
