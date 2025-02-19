package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/drive/v3"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db"
	"github.com/mmoehabb/studio-shop/handlers/guih"
	"github.com/mmoehabb/studio-shop/handlers/photo"
	"github.com/mmoehabb/studio-shop/handlers/section"
	"github.com/mmoehabb/studio-shop/handlers/user"
	"github.com/mmoehabb/studio-shop/middlewares"
	"github.com/mmoehabb/studio-shop/pages"
)

func main() {
	// initialize a context to share data between different templ components
	ctx := context.WithValue(context.Background(), "version", "v0.0.4")
	app := fiber.New()
	app.Static("/public", "./public/")

	// shall be used once and commented afterwards,
	// and maybe completed removed in production.
	app.Get("/seed", func(c *fiber.Ctx) error {
		defer anc.Recover(c)
		anc.Must(nil, db.Seed())

    service := anc.GetDriveService()
    files := anc.Must(service.Files.List().Do()).(*drive.FileList)
    
    // TODO: insert sections in Database
    // TODO: insert photos in Database

    // directory paterns:
    // 1. main-dir
    // 2. sec-dir
    // 1.1 inner-main-dir
    // 2.1 inner-sec--dir

    for _, file := range files.Files {
      log.Println(file.Name)
      log.Println("http://drive.google.com/uc?id=" + file.Id)
    }

		return c.SendString("Database has been seeded.")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Index().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Login().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Get("/gui/fragments/home", guih.HomeFragment)
	app.Get("/gui/fragments/my-cart", guih.MyCartFragment)
	app.Get("/gui/fragments/contact-us", guih.ContactUsFragment)
  app.Get("/gui/fragments/photo/:id", guih.PhotoFragment)

	app.Post("/login", user.Login)

  // ******** Auth Middleware ******** //
	app.Use(middlewares.Auth)

	app.Get("/admin", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Admin().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Get("/gui/fragments/dashboard", guih.DashboardFragment)
	app.Get("/gui/forms/add-section", guih.AddSectionForm)
	app.Get("/gui/forms/add-photo", guih.AddPhotoForm)

	app.Post("/section/add", section.Add)
  app.Delete("/section/delete/:id", section.Delete)

	app.Post("/photo/add", photo.Add)
  app.Delete("/photo/delete/:id", photo.Delete)

	app.Listen(":8080")
}
