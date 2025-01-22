package about

import "github.com/MickDuprez/gobase/internal/core/interfaces"

func New() interfaces.Feature {
	return interfaces.Feature{
		Name: "about",
		Path: "internal/features/about",
		NavItems: []interfaces.NavItem{
			{
				Title:    "About",
				URL:      "/about",
				Priority: 10,
			},
		},
		Routes: setupRoutes,
	}
}

func setupRoutes(app interfaces.App) {
	h := &Handler{app: app}
	app.Handle("GET /about", h.about)
	app.Handle("GET /about/team", h.team)
	app.Handle("GET /about/contact", h.contact)
}
