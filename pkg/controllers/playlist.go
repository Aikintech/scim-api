package controllers

import "github.com/gofiber/fiber/v2"

func CreatePlaylist(c *fiber.Ctx) error {
	return c.SendString("Ok")
}
