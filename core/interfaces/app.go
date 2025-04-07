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

	// Session helpers
	SessionSetValue(r *http.Request, key string, value interface{}) error
	SessionGetValue(r *http.Request, key string) (interface{}, bool)
	SessionGetString(r *http.Request, key string) (string, bool)
	SessionGetInt(r *http.Request, key string) (int64, bool)
}

type Feature struct {
	Name     string
	Path     string
	NavItems []NavItem
	Routes   func(app App)
	OnInit   func(app App) error // Optional initialization
}

type NavItem struct {
	Title     string
	URL       string // only used for sub items
	Priority  int
	SubItems  []NavItem // list of navigation links that use the URL field
	IsDivider bool      // optional marker for a divider
}
