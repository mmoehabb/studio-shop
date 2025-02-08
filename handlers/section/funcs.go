package section

import (
	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db/relations"
	"github.com/mmoehabb/studio-shop/db/sections"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx) error {
  defer anc.Recover(c)
  body := new(AddSectionBody)
  if err := c.BodyParser(body); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }

  ok, errs := ValidateAddSectionBody(body)
  if ok == false {
    return c.Status(fiber.StatusBadRequest).JSON(errs)
  }

  newSection := sections.DataModel { Title: body.Title }
  anc.Must(nil, sections.Add([]sections.DataModel{ newSection }))

  if body.Parent == "none" {
    return c.SendStatus(fiber.StatusOK)
  }

  newRelation := relations.DataModel { 
    Parent: anc.Must(strconv.Atoi(body.Parent)).(int),
    Child: anc.Must(sections.GetId(body.Title)).(int),
  }
  anc.Must(nil, relations.Add([]relations.DataModel{ newRelation }))

  return c.SendStatus(fiber.StatusOK)
}

func Delete(c *fiber.Ctx) error {
  defer anc.Recover(c)
  id, err := c.ParamsInt("id")
  if err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }
  anc.Must(nil, sections.Delete([]int{ id }))
  return c.SendStatus(fiber.StatusOK)
}
