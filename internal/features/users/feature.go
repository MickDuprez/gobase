package users

import "github.com/MickDuprez/gobase/internal/core/interfaces"

func New() interfaces.Feature {
	return interfaces.Feature{
		Name: "users",
		Path: "internal/features/users",
		NavItems: []interfaces.NavItem{
			{
				Title:    "Profile",
				URL:      "/profile",
				Priority: 90,
			},
			{
				Title:    "Login",
				URL:      "/login",
				Priority: 100,
			},
		},
		Routes: setupRoutes,
	}
}

func setupRoutes(app interfaces.App) {
	h := &Handler{app: app}

	// Auth routes
	app.Handle("GET /login", h.loginForm)
	app.Handle("POST /login", h.login)
	app.Handle("GET /register", h.registerForm)
	app.Handle("POST /register", h.register)
	app.Handle("POST /logout", h.logout)

	// Protected routes
	app.Handle("GET /profile", app.RequireAuth(h.profile))

	// htmx routes
	app.Handle("GET /profile/info/add", app.RequireAuth(h.addProfileInfo))
	app.Handle("POST /profile/info/save", app.RequireAuth(h.saveProfileInfo))
	app.Handle("GET /profile/info/show", app.RequireAuth(h.showProfileInfo))
}
