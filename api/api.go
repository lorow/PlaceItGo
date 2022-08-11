package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func serveIndex(c *fiber.Ctx) error{
	return c.SendString("Hello, World!")
}


func getImage(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("%s %s %s", c.Params("animal"), c.Params("width"), c.Params("height")))
}

func StartServer() error {
	app := fiber.New()
	app.Get("/",  serveIndex)
	app.Get("/:animal/:width/:height", getImage)
	return app.Listen(":8080")
}