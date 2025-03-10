package user

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/mmoehabb/studio-shop/db/users"
	"github.com/mmoehabb/studio-shop/ui/forms"
)

// login hanlder for fiber endpoint /login
// it expects a POST request
func Login(c *fiber.Ctx) error {
	creds := new(Credentials)
	if err := c.BodyParser(creds); err != nil {
		return err
	}

	ok, errs := ValidateCreds(creds)
	if ok == false {
		forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	firstUser := users.IsEmpty()
	if firstUser {
		users.Add(creds.Username, creds.Password)
		c.Cookie(&fiber.Cookie{Name: "username", Value: creds.Username})
		c.Cookie(&fiber.Cookie{Name: "password", Value: creds.Password})
		return c.Redirect("/admin")
	}

	res, err := users.Get(creds.Username)
	if err != nil {
		errs["username"] = err.Error()
		forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
		return c.SendStatus(fiber.StatusNotFound)
	}

	if creds.Password != res.Password {
		errs["password"] = "wrong password."
		forms.Login(errs).Render(context.Background(), c.Response().BodyWriter())
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{Name: "username", Value: creds.Username})
	c.Cookie(&fiber.Cookie{Name: "password", Value: creds.Password})

	return c.Redirect("/admin")
}
