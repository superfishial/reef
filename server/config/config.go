package config

import (
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
	S3     S3Config     `koanf:"s3"`
}

type ServerConfig struct {
	Port                   uint16 `koanf:"port"`
	RootURL                string `koanf:"rootURL"`
	JwtSecret              string `koanf:"jwtSecret"`
	CookieEncryptionSecret string `koanf:"cookieEncryptionSecret"`
	AuthCookieName         string `koanf:"authCookieName"`
	StateCookieName        string `koanf:"stateCookieName"`
}

type OAuthConfig struct {
	ClientID         string `koanf:"clientID"`
	ClientSecret     string `koanf:"clientSecret"`
	AuthEndpoint     string `koanf:"authEndpoint"`
	TokenEndpoint    string `koanf:"tokenEndpoint"`
	UserInfoEndpoint string `koanf:"userInfoEndpoint"`
	DeviceEndpoint   string `koanf:"deviceEndpoint"`
}

type S3Config struct {
	Bucket         string `koanf:"bucket"`
	Region         string `koanf:"region"`
	Endpoint       string `koanf:"endpoint"`
	AccessKey      string `koanf:"accessKey"`
	SecretKey      string `koanf:"secretKey"`
	ForcePathStyle bool   `koanf:"forcePathStyle"`
}

// TODO: Should error if can't values for required values like
// JwtSecret, CookieEncryptionSecret, etc.
func LoadConfig() (Config, error) {
	// Default -> Yaml -> Env
	k.Load(confmap.Provider(map[string]interface{}{
		"server": map[string]interface{}{
			"port":    3000,
			"rootURL": "http://localhost:3000",
		},
		"oauth": map[string]interface{}{},
		"s3": map[string]interface{}{
			"bucket":         "reef",
			"region":         "us-central1",
			"forcePathStyle": false,
		},
	}, "."), nil)
	k.Load(file.Provider("config.yaml"), yaml.Parser())
	k.Load(env.Provider("REEF_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "REEF_")), "_", ".", -1)
	}), nil)

	// Unmarshal into struct
	var conf Config
	err := k.Unmarshal("", &conf)
	return conf, err
}
