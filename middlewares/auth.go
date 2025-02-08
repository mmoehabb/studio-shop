package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mmoehabb/studio-shop/db/users"
)

func Auth(c *fiber.Ctx) error {
	if c.Cookies("username") == "" {
		return c.Redirect("/login")
	}
	data, err := users.Get(c.Cookies("username"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if c.Cookies("password") != data.Password {
		return c.Redirect("/login")
	}
	return c.Next()
}
