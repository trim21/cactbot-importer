package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"
	_ "time/tzdata" //内嵌时区数据资源

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"cactbot_importer/pkg/fetch"
	"cactbot_importer/pkg/repo"
)

func main() {
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
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	SetupRouter(app)

	logrus.Infoln("http server started")
	logrus.Fatalln(app.Listen(":3002"))
}

//go:embed static
var static embed.FS

func SetupRouter(app fiber.Router) {
	type Urls struct {
		Urls []string `json:"urls" form:"url" xml:"urls"`
	}
	app.Options("/v0/raidboss.js", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
		c.Set(fiber.HeaderAccessControlAllowHeaders, "content-type")
		return nil
	})
	app.Post("/v0/raidboss.js", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
		c.Set(fiber.HeaderAccessControlAllowHeaders, "content-type")
		p := new(Urls)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		logrus.WithField("urls", p.Urls).Infoln("post")

		js, err := fetch.Fetch(p.Urls)
		if err != nil {
			if errors.Is(err, fetch.ErrTransform) {
				return c.Status(502).SendString(err.Error() + "\n" + js)
			} else if errors.Is(err, repo.ErrNestedRepo) {
				return c.Status(502).SendString(err.Error())
			}
			fmt.Println(err)
			return err
		}

		return c.SendString(js)
	})

	app.Use(filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "static",
		Index:      "index.html",
	}))
}
