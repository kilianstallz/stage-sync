package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "",
		AppName:      "stage-sync-api",
	})
	app.Use(logger.New())
	app.Use(pprof.New())
	app.Use(recover2.New())
	app.Use(requestid.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	creds := app.Group("/creds")
	creds.Post("create", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Create creds")
	})

	log.Fatal(app.Listen(":8080"))

}
