package auth

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/superfishial/reef/server/config"
)

func MiddlewareHandler(conf config.ServerConfig, requireAuth bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Check if user is already authenticated
		headerOrCookie := c.Cookies(conf.AuthCookieName)
		if headerOrCookie == "" {
			re := regexp.MustCompile(`Bearer (.*)`)
			matches := re.FindStringSubmatch(c.Get("Authorization"))
			if len(matches) > 1 {
				headerOrCookie = matches[1]
			}
		}
		if headerOrCookie == "" {
			if requireAuth {
				return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("no token found in cookie '%s' or Authorization header in the form Bearer <token>", conf.AuthCookieName))
			}
			return c.Next()
		}

		// Verify JWT token
		token, err := jwt.Parse(headerOrCookie, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != "HS256" {
				return nil, errors.New("invalid signing algorithm")
			}
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
