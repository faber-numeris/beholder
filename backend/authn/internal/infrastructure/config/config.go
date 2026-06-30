package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type IAppConfig interface {
	IServiceConfig
	IDatabaseConfig
	IMailConfig
}

type AppConfig struct {
	ServiceConfig
	DatabaseConfig
	MailConfig
}

var configErr error

func NewConfig() (IAppConfig, error) {
	cfg, err := env.ParseAs[AppConfig]()
	if err != nil {
		configErr = fmt.Errorf("failed to parse environment variables: %w", err)
		return nil, err
	}
	return &cfg, configErr
}
