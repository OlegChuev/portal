package endpoints

import (
	"net/http"
	"portal/actions/render"
	"portal/models"

	"github.com/gobuffalo/buffalo"
)

// Index renders index page with all users.
func Index(c buffalo.Context) error {
	users := []models.User{}
	err := models.DB.All(&users)

	if err != nil {
		c.Set("users", []models.User{})
		c.Flash().Add("danger", "There is an error in fetching users")
	} else {
		c.Set("users", users)
	}

	template := render.GetRender().HTML("user/index.html", "layout/dashboard.html")

	return c.Render(http.StatusOK, template)
}
