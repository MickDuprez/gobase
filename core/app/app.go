package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MickDuprez/gobase/core/auth"
	"github.com/MickDuprez/gobase/core/config"
	"github.com/MickDuprez/gobase/core/database"
	"github.com/MickDuprez/gobase/core/interfaces"
	"github.com/MickDuprez/gobase/core/middleware"
	"github.com/MickDuprez/gobase/core/template"
)

type Application struct {
	templates      *template.Manager
	mux            *http.ServeMux
	features       map[string]*interfaces.Feature
	auth           *auth.AuthDB
	securityConfig *middleware.SecurityConfig
	db             *database.DB
}

func (app *Application) DB() *database.DB {
	return app.db
}

func New(cfg *config.AppConfig) (*Application, error) {
	// Initialize auth
	authDB, err := auth.NewAuthDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth: %w", err)
	}

	// Initialize template manager
	tm, err := template.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}

	// Initialize database - non-fatal
	var db *database.DB
	db, err = database.New(cfg.DBConfig)
	if err != nil {
		log.Printf("WARNING: Database initialization failed: %v", err)
		// Continue with nil db
	}

	app := &Application{
		templates:      tm,
		mux:            http.NewServeMux(),
		features:       make(map[string]*interfaces.Feature),
		auth:           authDB,
		db:             db, // Might be nil!
		securityConfig: cfg.SecConfig,
	}

	// Add static file server
	fileServer := http.FileServer(http.Dir("static"))
	app.mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	return app, nil
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sw := &statusWriter{ResponseWriter: w}

	app.mux.ServeHTTP(sw, r)

	log.Printf(
		"%s %s %d %v",
		r.Method,
		r.URL.Path,
		sw.status,
		time.Since(start),
	)
}

func (app *Application) Handle(pattern string, handler http.HandlerFunc) {
	secureHandler := middleware.SecurityHeaders(app.securityConfig)(handler)
	app.mux.HandleFunc(pattern, secureHandler)
}

func (app *Application) RenderTemplate(w http.ResponseWriter, r *http.Request, feature, page string, data interface{}) error {
	return app.templates.Render(w, r, feature, page, data)
}

func (app *Application) RenderPartial(w http.ResponseWriter, r *http.Request, feature, partial string, data interface{}) error {
	return app.templates.RenderPartial(w, r, feature, partial, data)
}

func (app *Application) RegisterHelperFunc(name string, fn interface{}) {
	app.templates.RegisterHelperFunc(name, fn)
}

func (app *Application) RegisterFeature(f interfaces.Feature) error {
	// Register feature's templates
	if err := app.templates.RegisterFeature(f.Name, f.Path, f.NavItems...); err != nil {
		return err
	}

	// Set up feature's routes
	f.Routes(app)

	app.features[f.Name] = &f
	return nil
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(b)
}

func (app *Application) Auth() *auth.AuthDB {
	return app.auth
}

func (app *Application) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return app.auth.RequireAuth(next)
}
