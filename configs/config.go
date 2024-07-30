package configs

import (
	"github.com/joho/godotenv"
	"log"
)

type IConfig interface {
	App() IAppConfig
}

type IAppConfig interface {
	GetToken() string
}

type config struct {
	app *app
}

type app struct {
	token string
}

func (c *config) App() IAppConfig {
	return c.app
}

func (a *app) GetToken() string {
	return a.token
}

func NewConfig(envpath string) IConfig {
	envMap, err := godotenv.Read(envpath)
	if err != nil {
		log.Fatal("failed to read config from .env")
	}

	return &config{
		app: &app{
			token: envMap["APP_TOKEN"],
		},
	}
}
