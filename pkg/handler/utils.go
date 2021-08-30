package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

func templateRender(c *fiber.Ctx, name string, bind interface{}) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if err := c.App().Config().Views.Render(buf, name, bind); err != nil {
		return err
	}

	return c.Send(buf.Bytes())
}
