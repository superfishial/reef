package auth

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/superfishial/reef/server/config"
)

func LoginHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		state, err := GenerateStateToken()
		if err != nil {
			return err
		}

		rootURL, err := url.Parse(conf.Server.RootURL)
		if err != nil {
			return fmt.Errorf("failed to parse rootURL from config: %w", err)
		}
		c.Cookie(&fiber.Cookie{
			Name:     conf.Server.StateCookieName,
			Value:    state,
			HTTPOnly: true,
			Secure:   rootURL.Scheme == "https", // TODO: Make true for production
		})

		authCodeURL := getOAuthClient(conf).AuthCodeURL(state)
		c.Redirect(authCodeURL)
		return nil
	}
}
