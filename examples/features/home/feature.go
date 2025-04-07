package home

import (
	"net/http"

	"github.com/MickDuprez/gobase/core/interfaces"
)

func New() interfaces.Feature {
	return interfaces.Feature{
		Name: "home",
		Path: "features/home",
		NavItems: []interfaces.NavItem{
			{
				Title:    "Main",
				Priority: 0,
				SubItems: []interfaces.NavItem{
					{Title: "Home", URL: "/"},
					{Title: "Overview", URL: "/overview"},
					{Title: "Stats", URL: "/stats"},
					{IsDivider: true},
					{Title: "Settings", URL: "/settings"},
				},
			},
		},
		Routes: setupRoutes,
	}
}

func setupRoutes(app interfaces.App) {
	h := &Handler{app: app}
	app.Handle("GET /", h.home)
}

// internal/features/home/handler.go
type Handler struct {
	app interfaces.App
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	_ = h.app.RenderTemplate(w, r, "home", "home", nil)
}
