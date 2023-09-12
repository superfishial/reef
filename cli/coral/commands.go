package coral

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/superfishial/reef/cli/config"
)

func GetCommand(conf config.Config) *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get remote file",
		Action: func(c *cli.Context) error {
			name := c.Args().First()
			if name == "" {
				return fmt.Errorf("no file name provided")
			}
			body, err := Get(conf.OAuth.Token, conf.Server.RootURL, name)
			if err != nil {
				return err
			}
			defer body.Close()
			io.Copy(os.Stdout, body)
			return nil
		},
	}
}
