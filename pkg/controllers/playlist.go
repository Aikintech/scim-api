package controllers

import "github.com/gofiber/fiber/v2"

func ClientCreatePlaylist(c *fiber.Ctx) error {
	return c.SendString("Ok")
}
