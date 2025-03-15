package photo

import (
	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db/photos"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx) error {
	defer anc.Recover(c)
	body := new(AddPhotoBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	ok, errs := ValidateAddPhotoBody(body)
	if ok == false {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	newPhoto := photos.DataModel{
		Name:      body.Name,
		Url:       body.Url,
		SectionId: anc.Must(strconv.Atoi(body.SectionId)).(int),
	}
	anc.Must(nil, photos.Add([]photos.DataModel{newPhoto}))

	return c.SendStatus(fiber.StatusOK)
}

func Delete(c *fiber.Ctx) error {
	defer anc.Recover(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	anc.Must(nil, photos.Delete(id))
	return c.SendStatus(fiber.StatusOK)
}
