package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/superfishial/reef/cli/auth"
	"github.com/superfishial/reef/cli/config"
)

func main() {
	log.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true, DisableTimestamp: true})

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	app := &cli.App{
		Name:  "reef",
		Usage: "flexible remote file storage",
		Commands: []*cli.Command{
			auth.GetCommand(conf),
			config.GetCommand(),
			{
				Name:  "get",
				Usage: "Get a file",
				Action: func(c *cli.Context) error {
					fmt.Println("getting file: ", c.Args().First())
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List files",
				Action: func(c *cli.Context) error {
					fmt.Println("putting file: ", c.Args().First())
					return nil
				},
			},
			{
				Name:  "put",
				Usage: "Put a file",
				Action: func(c *cli.Context) error {
					fmt.Println("putting file: ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del"},
				Usage:   "Delete a file",
				Action: func(c *cli.Context) error {
					fmt.Println("deleting file: ", c.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
