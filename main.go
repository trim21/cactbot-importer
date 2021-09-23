package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"cactbot_importer/pkg/handler"
)

//go:generate pkger
func main() {
	app := fiber.New(fiber.Config{
		StrictRouting:         true,
		CaseSensitive:         true,
		DisableStartupMessage: true,
	})
	app.Use(recover.New())

	handler.SetupRouter(app)

	fmt.Println("http://127.0.0.1:5000/")
	log.Fatalln(app.Listen(":3002"))
}
