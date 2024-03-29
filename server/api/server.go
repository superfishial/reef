package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"

	"github.com/superfishial/reef/server/auth"
	"github.com/superfishial/reef/server/config"
	"github.com/superfishial/reef/server/coral"
)

func StartServer(conf config.Config) {
	app := fiber.New(fiber.Config{
		AppName:           "reef",
		BodyLimit:         10 * 1024 * 1024 * 1024,
		StreamRequestBody: true,
	})

	app.Use(logger.New(logger.Config{
		Format: "${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: conf.Server.CookieEncryptionSecret,
	}))

	app.Route("/v1", func(v1 fiber.Router) {
		// Authentication routes
		v1.Route("/auth", func(authR fiber.Router) {
			authR.Get("/login", auth.LoginHandler(conf))
			authR.Get("/callback", auth.CallbackHandler(conf))
			authR.Get("/sign", auth.SignHandler(conf.Server))
		}, "auth.")

		// Config routes
		v1.Get("/config/cli", config.CLIConfigHandler(conf))

		// Coral (file) routes
		v1.Delete("/coral/:filename", auth.MiddlewareHandler(conf.Server, true), func(c *fiber.Ctx) error {
			err := coral.DeleteFile(conf.S3, c.Params("filename"))
			if err != nil {
				return err
			}
			return c.SendStatus(204)
		})

		v1.Put("/coral/:filename", auth.MiddlewareHandler(conf.Server, true), func(c *fiber.Ctx) error {
			tokenClaims := c.Locals("token").(*jwt.Token)
			ownerSub, err := tokenClaims.Claims.GetSubject()
			if err != nil {
				return fiber.ErrInternalServerError
			}
			metadata := coral.FileMetadata{OwnerSub: ownerSub, Public: c.QueryBool("public", false)}

			err = coral.UploadFile(
				conf.S3,
				c.Params("filename"),
				metadata,
				c.Context().RequestBodyStream(),
			)
			if err != nil {
				return errors.Join(fiber.ErrInternalServerError, err)
			}
			return c.SendStatus(201)
		})

		v1.Get("/coral/:filename", auth.MiddlewareHandler(conf.Server, false), func(c *fiber.Ctx) error {
			_, err := coral.GetFileMetadata(conf.S3, c.Params("filename"))
			// TODO: Handle 404s
			// TODO: Check if the user has access to read file
			if err != nil {
				return errors.Join(fiber.ErrInternalServerError, err)
			}

			stream, err := coral.DownloadFile(conf.S3, c.Params("filename"))
			if err != nil {
				return errors.Join(fiber.ErrInternalServerError, err)
			}
			return c.SendStream(stream)
		})
	}, "v1.")

	log.Fatal(app.Listen(":3000"))
}
