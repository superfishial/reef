package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/superfishial/reef/server/config"
)

func MiddlewareHandler(conf config.ServerConfig, requireCookie bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Check if user is already authenticated
		cookie := c.Cookies(conf.AuthCookieName)
		if cookie == "" {
			if requireCookie {
				return fiber.NewError(fiber.StatusUnauthorized, "no token found in cookie")
			}
			return nil
		}

		// Verify JWT token
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JwtSecret), nil
		})
		if err != nil {
			return errors.Join(
				fiber.NewError(fiber.StatusInternalServerError, "failed while parsing token from cookie"),
				err,
			)
		}
		if !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token from cookie")
		}

		// Store the token so we can access it in other handlers
		c.Locals("token", token)

		return c.Next()
	}
}
