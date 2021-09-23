package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"cactbot_importer/pkg/handler"
)

//go:generate pkger
func main() {
	setTimeZone()
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339,
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		PrettyPrint:       false,
	})
	app := fiber.New(fiber.Config{
		StrictRouting:         true,
		CaseSensitive:         true,
		DisableStartupMessage: true,
	})
	app.Use(recover.New())

	handler.SetupRouter(app)

	logrus.Infoln("http server started")
	logrus.Fatalln(app.Listen(":3002"))
}

func setTimeZone() {
	time.Local = time.FixedZone("Asia/Shanghai", 8*3600)
}
