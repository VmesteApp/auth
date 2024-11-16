package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP             `yaml:"http"`
		GRPC             `yaml:"grpc"`
		Log              `yaml:"logger"`
		PG               `yaml:"postgres"`
		VkAPI            `yaml:"vk_api"`
		JwtConfig        `yaml:"jwt"`
		SuperAdminConfig `yaml:"superadmin"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"pg_url" env:"PG_URL"`
	}

	VkAPI struct {
		AppId      int    `env-required:"true" yaml:"app_id" env:"VK_APP_ID"`
		PrivateKey string `env-required:"true" env:"VK_PRIVATE_KEY"`
		ServiceKey string `env-required:"true" env:"VK_SERVICE_KEY"`
	}

	JwtConfig struct {
		Secret string        `env-required:"true" env:"JWT_TOKEN_SECRET"`
		TTL    time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
	}

	SuperAdminConfig struct {
		Email    string `env-required:"true" env:"SUPER_ADMIN_EMAIL"`
		Password string `env-required:"true" env:"SUPER_ADMIN_PASSWORD"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read yml config: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	return cfg, nil
}
