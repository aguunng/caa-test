package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog/log"
)

func Load() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Msgf("unable to parse env: %s", err.Error())
	}

	return &c
}

func ReadConfig() (*AppConfig, error) {
    file, err := os.Open(configFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config AppConfig
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func WriteConfig(config *AppConfig) error {
    file, err := os.Create(configFile)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(config)
}

type Config struct {
	App      App
	Qiscus   Qiscus
}

type App struct {
	SecretKey string `env:"APP_SECRET_KEY"`
}

type Qiscus struct {
	AppID       string `env:"QISCUS_APP_ID"`
	SecretKey   string `env:"QISCUS_SECRET_KEY"`
	Omnichannel Omnichannel
}

type Omnichannel struct {
	URL string `env:"QISCUS_OMNICHANNEL_URL"`
}

const configFile = "config.json"

type AppConfig struct {
	MaxCustomer int `json:"max_customer"`
}

