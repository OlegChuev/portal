package routes

import (
	"portal/actions/assets"
	"portal/actions/endpoints"

	"github.com/gobuffalo/buffalo"
)

func SetRoutes(app *buffalo.App) {
	setAuthRoutes(app)

	app.ServeFiles("/", assets.GetAssetsBox()) // serve files from the public directory
}

func setAuthRoutes(app *buffalo.App) {
	// Login
	auth := app.Group("/auth")
	auth.GET("/login", endpoints.LogInPage)
	auth.POST("/login", endpoints.LogIn)
	auth.GET("/logout", endpoints.LogOut)

	// Registration
	registration := app.Group("/registration")
	registration.GET("/sign_up", endpoints.SignUpPage)
	registration.POST("/sign_up", endpoints.SignUp)
}
