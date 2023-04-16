package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/superfishial/reef/server/config"
)

func CallbackHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Verify that we got the expected state to prevent CSRF attacks
		state := c.Query("state")
		expectedState := c.Cookies(conf.Server.StateCookieName)
		if state != expectedState {
			return fiber.NewError(fiber.StatusUnauthorized, "expected state did not equal received state")
		}

		// Receive access token from OAuth provider
		var code = c.Query("code")
		token, err := getOAuthClient(conf).Exchange(context.Background(), code)
		if err != nil {
			return errors.Join(
				fiber.NewError(fiber.StatusBadRequest, "failed to exchange code for token"),
				err,
			)
		}

		// Get user info with access token
		userInfo, err := GetUserInfo(token.AccessToken)
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, err)
		}

		// Sign token
		tokenString, err := SignToken(conf.Server.JwtSecret, userInfo)
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, err)
		}

		// Write access token to cookie
		c.Cookie(&fiber.Cookie{
			Name:    conf.Server.AuthCookieName,
			Value:   tokenString,
			Expires: time.Now().Add(90 * 24 * time.Hour),
		})

		return c.Redirect("/")
	}

}
