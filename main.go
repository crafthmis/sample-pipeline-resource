package main

import (
	"encoding/json"
	"net/http"
	"techbridge/pipeline/files"
	"techbridge/pipeline/models"
	"techbridge/pipeline/resources"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/resource", getDeploymentFile)
	app.Listen(":8080")
}

func getDeploymentFile(c *fiber.Ctx) (err error) {
	var payload models.Payload
	resource := c.Query("type", "deployment")

	if err := json.Unmarshal([]byte(c.Body()), &payload); err != nil {
		c.Status(http.StatusBadRequest).Send([]byte("Invalid request payload"))
		return err
	}

	file, err := resources.GetResource(payload, resource)

	if err != nil {
		c.Status(http.StatusInternalServerError).Send([]byte("Could not generate deployment file"))
		return
	}

	defer files.DeleteFile(file)
	return c.Type("application/octet-stream").Status(http.StatusCreated).SendFile(file, true)
}
