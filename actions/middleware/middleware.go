package middleware

import (
	"net/http"
	"portal/actions/endpoints"
	"portal/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func SetMiddleware(app *buffalo.App) {
	app.Use(SetCurrentUser)
	app.Use(Authorize)

	app.Middleware.Skip(Authorize, endpoints.SignUpPage, endpoints.SignUp, endpoints.LogIn, endpoints.LogInPage)
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)

			if err != nil {
				c.Logger().Warnf("user attempted to access with current_user_id '%v' that is not found: %v", uid, err)

				c.Session().Delete("current_user_id")
				c.Session().Set("redirectURL", c.Request().URL.String())
				c.Flash().Add("danger", "You must be authorized with a correct user to see that page")

				return c.Redirect(http.StatusFound, "/auth/login")
			}
			c.Set("current_user", u)
		}

		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")

			return c.Redirect(http.StatusFound, "/auth/login")
		}

		return next(c)
	}
}
