package interfaces

import (
	"net/http"

	"github.com/MickDuprez/gobase/core/auth"
	"github.com/MickDuprez/gobase/core/database"
)

type App interface {
	Handle(pattern string, handler http.HandlerFunc)
	RenderTemplate(w http.ResponseWriter, r *http.Request, feature, page string, data interface{}) error
	RenderPartial(w http.ResponseWriter, r *http.Request, feature, partial string, data interface{}) error
	RegisterFeature(f Feature) error
	Auth() *auth.AuthDB
	RequireAuth(next http.HandlerFunc) http.HandlerFunc
	DB() *database.DB
}

type Feature struct {
	Name     string
	Path     string
	NavItems []NavItem
	Routes   func(app App)
	OnInit   func(app App) error // Optional initialization
}

type NavItem struct {
	Title    string
	URL      string
	Priority int
}
