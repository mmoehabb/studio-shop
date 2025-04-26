package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/handlers/database"
	"github.com/mmoehabb/studio-shop/handlers/guih"
	"github.com/mmoehabb/studio-shop/handlers/photo"
	"github.com/mmoehabb/studio-shop/handlers/section"
	"github.com/mmoehabb/studio-shop/handlers/user"
	"github.com/mmoehabb/studio-shop/middlewares"
	"github.com/mmoehabb/studio-shop/pages"
)

func main() {
	ancillaries.LoadEnv()
	// initialize a context to share data between different templ components
	ctx := context.WithValue(context.Background(), "version", "v1.0.0")
	app := fiber.New()
	app.Static("/public", "./public/")
	app.Static("/apk", "./public/rashed-studio.apk")

	app.Use(cache.New(cache.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			key := struct {
				Path    string
				Body    string
				Queries map[string]string
			}{
				Path:    c.Path(),
				Body:    string(c.Body()),
				Queries: c.Queries(),
			}

			k, err := json.Marshal(key)
			if err != nil {
				utils.CopyString(c.Path())
			}
			return utils.CopyString(string(k))
		},
	}))

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
	app.Get("/gui/fragments/my-cart/list", guih.MyCartList)
	app.Get("/gui/fragments/contact-us", guih.ContactUsFragment)
	app.Get("/gui/fragments/photo/:id", guih.PhotoFragment)

	app.Get("/gui/lists/photos", guih.PhotosList)
	app.Get("/gui/empty", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/login", user.Login)

	// ******** Auth Middleware ******** //
	app.Use(middlewares.Auth)

	app.Get("/gui/fragments/dashboard", guih.DashboardFragment)
	app.Get("/gui/forms/add-section", guih.AddSectionForm)
	app.Get("/gui/forms/add-photo", guih.AddPhotoForm)

	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Add(fiber.HeaderCacheControl, "no-store")
		return c.Next()
	})

	app.Get("/admin", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		pages.Admin().Render(ctx, c.Response().BodyWriter())
		return c.SendStatus(200)
	})

	app.Get("/seed", database.Seed)
	app.Get("/reseed", database.Reseed)
	app.Get("/reseed/status", database.ReseedStatus)

	app.Post("/section/add", section.Add)
	app.Delete("/section/delete/:id", section.Delete)

	app.Post("/photo/add", photo.Add)
	app.Delete("/photo/delete/:id", photo.Delete)

	app.Listen(":8080")
}
