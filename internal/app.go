package internal

import (
	"errors"

	"github.com/spf13/viper"
)

type App struct {
	HomeDir       string
	GitDir        string
	Repo          string
	DefaultEditor string
	Config        *viper.Viper
}

func (app *App) SetConfigurationValue(key, value string) error {
	switch key {
	case "repository":
		app.Config.Set("repository", value)
	default:
		return errors.New("invalid configuration setting")
	}
	return app.Config.WriteConfig()
}
