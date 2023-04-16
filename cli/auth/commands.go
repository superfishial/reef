package auth

import (
	"github.com/superfishial/reef/cli/config"
	"github.com/urfave/cli/v2"
)

func GetCommand(conf config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "handle authentication",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login via oauth device flow",
				Action: func(c *cli.Context) error {
					token, err := performLoginFlow(conf)
					if err != nil {
						return err
					}
					return config.SetToken(token)
				},
			},
			{
				Name:  "logout",
				Usage: "logout and delete stored token",
				Action: func(c *cli.Context) error {
					return config.SetToken("")
				},
			},
		},
	}
}
