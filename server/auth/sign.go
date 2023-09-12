package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/superfishial/reef/server/config"
)

func SignHandler(conf config.ServerConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessToken := c.Query("access_token")
		if accessToken == "" {
			return fiber.NewError(fiber.StatusBadRequest, "access_token must be passed as a query param")
		}
		userInfo, err := GetUserInfo(accessToken)
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, err)
		}

		token, err := SignToken(conf.JwtSecret, userInfo)
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, err)
		}

		return c.SendString(token)
	}
}
