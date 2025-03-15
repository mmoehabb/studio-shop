package guih

import (
	"context"
	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db/sections"
	"github.com/mmoehabb/studio-shop/ui/forms"

	"github.com/gofiber/fiber/v2"
)

func AddSectionForm(c *fiber.Ctx) error {
	defer anc.Recover(c)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	notAlbumSections := anc.Must(sections.GetNotAlbums()).([]sections.DataModel)
	forms.AddSection(notAlbumSections, nil).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
	return c.SendStatus(fiber.StatusOK)
}

func AddPhotoForm(c *fiber.Ctx) error {
	defer anc.Recover(c)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	albumSections := anc.Must(sections.GetAlbums()).([]sections.DataModel)
	forms.AddPhoto(albumSections, nil).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
	return c.SendStatus(fiber.StatusOK)
}
