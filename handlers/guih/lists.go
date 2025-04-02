package guih

import (
	"context"
	"strconv"
	"strings"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db/photos"
	"github.com/mmoehabb/studio-shop/db/relations"
	"github.com/mmoehabb/studio-shop/ui/components"
	"github.com/mmoehabb/studio-shop/ui/fragments"

	"github.com/gofiber/fiber/v2"
)

func PhotosList(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	sectionId := c.QueryInt("section-id", 0)
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 6)

	isAlbumSection := relations.IsAlbum(sectionId)
	if isAlbumSection == false {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	list := anc.Must(photos.GetOfWithPagination(sectionId, page, size)).([]photos.DataModel)
	components.PhotoList(list, page, size).Render(context.Background(), c.Response().BodyWriter())

	return c.SendStatus(fiber.StatusOK)
}

func MyCartList(c *fiber.Ctx) error {
	itemsQuery := c.Query("items")
	if itemsQuery == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	itemsQuery = itemsQuery[1 : len(itemsQuery)-1]
	items := strings.Split(itemsQuery, ",")

	var ids = []int{}
	for _, item := range items {
		id, err := strconv.Atoi(item)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	photosList, _ := anc.Must(photos.GetList(ids)).([]photos.DataModel)
	fragments.MyCartList(photosList).Render(context.Background(), c.Response().BodyWriter())
	return c.SendStatus(fiber.StatusOK)
}
