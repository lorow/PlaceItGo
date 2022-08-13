package pkg

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func serveIndex(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func getImage(c *fiber.Ctx) error {

	animal := c.Params("animal")
	width_str := c.Params("width")
	height_str := c.Params("height")

	width, width_err := strconv.Atoi(width_str)
	height, height_err := strconv.Atoi(height_str)

	if width_err != nil || height_err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	imageService := ImageManager{}

	imageData, err := imageService.GetImage(animal, width, height)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "image/jpeg")
	return c.Send(imageData)
}

func StartServer() error {
	app := fiber.New()
	app.Get("/", serveIndex)
	app.Get("/:animal/:width/:height", getImage)
	return app.Listen(":8080")
}
