package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
)

func GetCommand() *cli.Command {
	return &cli.Command{
		Name: "config",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Usage:    "url of the reef server",
						Required: true,
					},
				},
				Action: initCommand,
			},
			{
				Name:   "print",
				Usage:  "Print configuration",
				Action: printCommand,
			},
		},
	}
}

func initCommand(c *cli.Context) error {
	serverUrlStr := c.String("server-url")
	serverUrl, err := url.Parse(serverUrlStr)
	if err != nil {
		return fmt.Errorf("invalid server url: %w", err)
	}
	serverUrl.Path = "/v1/config/cli"

	// Fetch the config
	resp, err := http.Get(serverUrl.String())
	if err != nil {
		return fmt.Errorf("failed to get config from server: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get config from server with status: %s", resp.Status)
	}
	config, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read config body from server: %w", err)
	}

	// Parse the config
	k = koanf.New(".")
	k.Load(rawbytes.Provider(config), yaml.Parser())

	// Merge and save the config
	err = MergeConfig(k)
	if err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}
	err = SaveConfig()
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}

func printCommand(c *cli.Context) error {
	yamlBytes, err := GetKoanf().Marshal(yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println(string(yamlBytes))
	return nil
}
