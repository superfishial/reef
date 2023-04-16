package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	Server ServerConfig `koanf:"server"`
	OAuth  OAuthConfig  `koanf:"oauth"`
}

type ServerConfig struct {
	RootURL string `koanf:"rootURL"`
}

type OAuthConfig struct {
	ClientID       string `koanf:"clientID"`
	TokenEndpoint  string `koanf:"tokenEndpoint"`
	DeviceEndpoint string `koanf:"deviceEndpoint"`
	Token          string `koanf:"token"`
}

func ConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return configDir + "/reef/cli.yaml"
}

func LoadConfig() (Config, error) {
	// Default -> Yaml -> Env
	k.Load(confmap.Provider(map[string]interface{}{
		"server": map[string]interface{}{
			"rootURL": "http://localhost:3000",
		},
		"oauth": map[string]interface{}{},
	}, "."), nil)
	k.Load(file.Provider(ConfigPath()), yaml.Parser())
	k.Load(env.Provider("REEF_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "REEF_")), "_", ".", -1)
	}), nil)

	// Unmarshal into struct
	var config Config
	err := k.Unmarshal("", &config)
	return config, err
}

func MergeConfig(newK *koanf.Koanf) error {
	return k.Merge(newK)
}

func GetKoanf() *koanf.Koanf {
	return k
}

func SaveConfig() error {
	yamlBytes, err := k.Marshal(yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.MkdirAll(path.Dir(ConfigPath()), 0700)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	err = os.WriteFile(ConfigPath(), yamlBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func SetToken(token string) error {
	k.Set("oauth.token", token)
	return SaveConfig()
}
