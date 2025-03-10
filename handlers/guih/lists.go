package guih

import (
	"context"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db/photos"
	"github.com/mmoehabb/studio-shop/db/relations"
	"github.com/mmoehabb/studio-shop/ui/components"

	"github.com/gofiber/fiber/v2"
)

func PhotosList(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  sectionId := c.QueryInt("section-id", 0)
  page := c.QueryInt("page", 1)
  size := c.QueryInt("size", 10)

  isAlbumSection := relations.IsAlbum(sectionId)
  if isAlbumSection == false {
    return c.SendStatus(fiber.StatusBadRequest)
  }

  list := anc.Must(photos.GetOfWithPagination(sectionId, page, size)).([]photos.DataModel)
  components.PhotoList(list).Render(context.Background(), c.Response().BodyWriter())

  return c.SendStatus(fiber.StatusOK)
}
