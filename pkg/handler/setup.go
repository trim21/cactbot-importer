package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"

	"cactbot_importer/pkg/fetch"
	"cactbot_importer/pkg/repo"
	"cactbot_importer/pkg/s"
	"cactbot_importer/pkg/utils"
)

func SetupRouter(app fiber.Router) {
	static := pkger.Dir("/static")

	app.Get("/js", func(c *fiber.Ctx) error {
		h, err := utils.GenerateSecureToken(64)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"title":  h,
			"length": len(h),
		})
	})

	app.Get("/js/:hash", func(c *fiber.Ctx) error {
		h := c.Params("hash")
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJavaScriptCharsetUTF8)
		return templateRender(c, "loader.js", s.Loader{UniqueID: h})
	})

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
