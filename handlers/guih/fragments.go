package guih

import (
	"context"

	anc "goweb/ancillaries"
	"goweb/db/photos"
	"goweb/db/relations"
	"goweb/db/sections"
	"goweb/ui/fragments"

	"github.com/gofiber/fiber/v2"
)

func DashboardFragment(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  sectionId := c.QueryInt("section-id", 0)

  if sectionId == 0 {
    list := anc.Must(sections.GetMain()).([]sections.DataModel)
    fragments.Dashboard(list).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusOK)
  } 

  isAlbumSection := relations.IsAlbum(sectionId)
  if isAlbumSection {
    list := anc.Must(photos.GetOf(sectionId)).([]photos.DataModel)
    fragments.PhotosDashboard(list).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusOK)
  }

  ids := anc.Must(relations.GetSectionsOf(sectionId)).([]int)
  list := anc.Must(sections.Get(ids)).([]sections.DataModel)
  fragments.Dashboard(list).Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}

func HomeFragment(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  sectionId := c.QueryInt("section-id", 0)

  if sectionId == 0 {
    list := anc.Must(sections.GetMain()).([]sections.DataModel)
    fragments.Home(list).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusOK)
  } 

  isAlbumSection := relations.IsAlbum(sectionId)
  if isAlbumSection {
    list := anc.Must(photos.GetOf(sectionId)).([]photos.DataModel)
    fragments.PhotosHome(list).Render(context.Background(), c.Response().BodyWriter())
    return c.SendStatus(fiber.StatusOK)
  }

  ids := anc.Must(relations.GetSectionsOf(sectionId)).([]int)
  list := anc.Must(sections.Get(ids)).([]sections.DataModel)
  fragments.Home(list).Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}

func MyCartFragment(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  fragments.MyCart().Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}

func ContactUsFragment(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  fragments.ContactUs().Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}

func PhotoFragment(c *fiber.Ctx) error {
  c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
  photoId := anc.Must(c.ParamsInt("id")).(int)

  photo, err := photos.Get(photoId)
  if err != nil {
    return c.SendStatus(fiber.StatusNotFound)
  }
  
  fragments.Photo(&photo).Render(context.Background(), c.Response().BodyWriter())
  return c.SendStatus(fiber.StatusOK)
}
