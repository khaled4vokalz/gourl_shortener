package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	"gopkg.in/yaml.v3"
)

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
	var config Config
	filePath := getPath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return config, nil
}
