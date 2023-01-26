package localconfig

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Frontend struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Path string `mapstructure:"path"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Channel struct {
	// Subscribe string `mapstructure:"subscribe"`
	Name string `mapstructure:"name"`
}

type Config struct {
	Frontend Frontend `mapstructure:"frontend"`
	Redis    Redis    `mapstructure:"redis"`
	Channel  Channel  `mapstructure:"channel"`
}

// LoadConfig reads the file from path and return Secret
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromBytes(data)
}

// LoadConfigFromBytes reads the secret file from data bytes
func LoadConfigFromBytes(data []byte) (*Config, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix("NOTIF")
	fang.SetConfigType("yaml")

	if err := fang.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	fang.Get("name")

	var cfg Config
	err := fang.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error loading creds: %v", err)
	}

	return &cfg, nil
}
