package endpoints

import (
	"net/http"
	"portal/actions/render"

	"github.com/gobuffalo/buffalo"
)

// DashboardPage renders dashboard main page.
func DashboardPage(c buffalo.Context) error {
	template := render.GetRender().HTML("dashboard/index.html", "layout/dashboard.html")

	return c.Render(http.StatusOK, template)
}
