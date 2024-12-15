package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const ENV_PREFIX = "GOURLAPP_"

func getPath() string {
	common.LoadEnv()
	config_path := os.Getenv("CONFIG_PATH")
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}
	if config_path == "" {
		config_path = fmt.Sprintf("configuration/%s.yaml", strings.ToLower(env))
	}
	return config_path
}

func LoadConfig() (Config, error) {
	k := koanf.New(".")
	var config Config
	filePath := getPath()
	if err := k.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		fmt.Printf("Error loading YAML file: %v\n", err)
		return Config{}, err
	}

	k.Load(env.Provider(ENV_PREFIX, ".", func(s string) string {
		return strings.Replace(
			strings.TrimPrefix(s, ENV_PREFIX), "_", ".", -1)
	}), nil)

	if err := k.Unmarshal("", &config); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return Config{}, err
	}

	return config, nil
}
