package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App   `yaml:"app"`
		Log   `yaml:"logger"`
		Redis `yaml:"redis"`
		RPC   `yaml:"rpc"`
		PG    `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Redis struct {
		Expiration int `env-required:"false" yaml:"expiration"   env:"REDIS_EXP"`
		Port       int `env-required:"true"  yaml:"port"         env:"REDIS_PORT"`
	}

	RPC struct {
		Port int `env-required:"true" yaml:"port"   env:"RPC_PORT"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
