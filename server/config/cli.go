package config

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	cliConfig "github.com/superfishial/reef/cli/config"
)

func CLIConfigHandler(conf Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cliConf, err := cliConfig.LoadConfig()
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, fmt.Errorf("cli config file not found: %s", err))
		}
		cliConf.Server.RootURL = conf.Server.RootURL
		cliConf.OAuth.ClientID = conf.OAuth.ClientID
		cliConf.OAuth.TokenEndpoint = conf.OAuth.TokenEndpoint
		cliConf.OAuth.DeviceEndpoint = conf.OAuth.DeviceEndpoint

		k = koanf.New(".")
		k.Load(structs.Provider(cliConf, "koanf"), nil)

		yamlBytes, err := k.Marshal(yaml.Parser())
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, fmt.Errorf("failed to marshal cli config to yaml: %s", err))
		}

		return c.Send(yamlBytes)
	}
}
