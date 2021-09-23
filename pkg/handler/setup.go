package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"cactbot_importer/pkg/fetch"
	"cactbot_importer/pkg/repo"
)

func SetupRouter(app fiber.Router) {
	static := pkger.Dir("/static")

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
		
		logrus.WithField("urls", p.Urls).Infoln("")

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
		Root:  static,
		Index: "index.html",
	}))

}
