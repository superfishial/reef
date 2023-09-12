package coral

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/superfishial/reef/server/config"
)

func DeleteHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		err := DeleteFile(conf.S3, c.Params("filename"))
		if err != nil {
			return err
		}
		return c.SendStatus(204)
	}
}

func PutHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenClaims := c.Locals("token").(*jwt.Token)
		ownerSub, err := tokenClaims.Claims.GetSubject()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		metadata := FileMetadata{OwnerSub: ownerSub, Public: c.QueryBool("public", false)}

		err = UploadFile(
			conf.S3,
			c.Params("filename"),
			metadata,
			c.Context().RequestBodyStream(),
		)
		if err != nil {
			return errors.Join(fiber.ErrInternalServerError, err)
		}
		return c.SendStatus(201)
	}
}

func GetHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		metadata, err := GetFileMetadata(conf.S3, c.Params("filename"))
		if err != nil {
			return fiber.ErrInternalServerError
		}
		fmt.Println(metadata, err)

		stream, err := DownloadFile(conf.S3, c.Params("filename"))
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.SendStream(stream)
	}
}

func ListHandler(conf config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		metadata, err := GetFileMetadata(conf.S3, c.Params("filename"))
		if err != nil {
			return fiber.ErrInternalServerError
		}
		fmt.Println(metadata, err)

		stream, err := DownloadFile(conf.S3, c.Params("filename"))
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.SendStream(stream)
	}
}
