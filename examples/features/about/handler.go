package about

import (
	"net/http"

	"github.com/MickDuprez/gobase/core/interfaces"
)

type Handler struct {
	app interfaces.App
}

func (h *Handler) about(w http.ResponseWriter, r *http.Request) {
	h.app.RenderTemplate(w, r, "about", "page", nil)
}

func (h *Handler) team(w http.ResponseWriter, r *http.Request) {
	team := []struct {
		Name string
		Role string
	}{
		{"John Doe", "Lead Developer"},
		{"Jane Smith", "Designer"},
	}
	h.app.RenderTemplate(w, r, "about", "team", team)
}

func (h *Handler) contact(w http.ResponseWriter, r *http.Request) {
	h.app.RenderTemplate(w, r, "about", "contact", nil)
}
