package config

import (
	"os"
	"path"

	dto "github.com/mojganchakeri/whatsapp-manager/internal/DTO"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

type Config interface {
	GetConfig() *dto.Config
}

type gookitConfig struct {
	cfg *dto.Config
}

func New() Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := path.Join(pwd, "./config.yaml")

	c := new(dto.Config)

	config.WithOptions(func(o *config.Options) {
		config.WithOptions(config.ParseDefault)
		config.WithOptions(config.ParseEnv)
		o.DecoderConfig.TagName = "config"
	})

	config.AddDriver(yamlv3.Driver)

	err = config.LoadFiles(configPath)
	if err != nil {
		panic(err)
	}

	err = config.BindStruct("", c)
	if err != nil {
		panic(err)
	}

	return &gookitConfig{
		cfg: c,
	}
}

func (c gookitConfig) GetConfig() *dto.Config {
	return c.cfg
}
