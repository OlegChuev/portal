package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"
	"portal/actions/render"
	"portal/models"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// SignUpPage renders sign up page.
func SignUpPage(c buffalo.Context) error {
	c.Set("user", models.User{})
	template := render.GetRender().HTML("auth/sign_up.html", "layout/auth.html")

	return c.Render(http.StatusOK, template)
}

// SignUp tries to sign up user in application.
func SignUp(c buffalo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	verrs, err := user.Create(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", user)
		c.Set("errors", verrs)
		template := render.GetRender().HTML("auth/sign_up.html", "layout/auth.html")

		return c.Render(http.StatusOK, template)
	}

	c.Session().Set("current_user_id", user.ID)
	c.Flash().Add("success", "Welcome to portal!")

	return c.Redirect(http.StatusFound, "/")
}

// LogInPage renders login page.
func LogInPage(c buffalo.Context) error {
	c.Set("user", models.User{})
	template := render.GetRender().HTML("auth/login.html", "layout/auth.html")

	return c.Render(http.StatusOK, template)
}

// LogIn tries to authorize user in application.
func LogIn(c buffalo.Context) error {
	fmt.Println("IAM HERE")
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(user.Email))).First(user)

	// helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")

		c.Set("errors", verrs)
		c.Set("user", user)

		template := render.GetRender().HTML("auth/login.html", "layout/auth.html")

		return c.Render(http.StatusUnauthorized, template)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// couldn't find an user with the supplied email address.
			return bad()
		}

		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user.Password))
	if err != nil {
		return bad()
	}

	c.Session().Set("current_user_id", user.ID)
	c.Flash().Add("success", "Welcome Back to Buffalo!")

	redirectURL := "/"

	if redir, ok := c.Session().Get("redirectURL").(string); ok && redir != "" {
		redirectURL = redir
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

// LogOut clears user's session.
func LogOut(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")

	return c.Redirect(http.StatusFound, "/login")
}
