package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MickDuprez/gobase/core/interfaces"
)

type Manager struct {
	pageTemplates    map[string]*template.Template
	partialTemplates map[string]*template.Template
	navItems         []interfaces.NavItem
}

func New() (*Manager, error) {
	return &Manager{
		pageTemplates:    make(map[string]*template.Template),
		partialTemplates: make(map[string]*template.Template),
		navItems:         make([]interfaces.NavItem, 0),
	}, nil
}

func (m *Manager) RegisterFeature(name, path string, navItems ...interfaces.NavItem) error {
	// Get all page templates for this feature
	pages, err := filepath.Glob(filepath.Join(path, "templates", "*.html"))
	if err != nil {
		return fmt.Errorf("failed to find feature templates: %w", err)
	}

	log.Printf("Registering feature %s:", name)
	for _, page := range pages {
		pageName := filepath.Base(page)
		if pageName == "layout.html" {
			continue // Skip layout file as it's handled separately
		}
		log.Printf("  Processing: %s", pageName)

		// Start with base template
		ts, err := template.ParseFiles("templates/layouts/base.html")
		if err != nil {
			return fmt.Errorf("failed to parse base template: %w", err)
		}

		// Add feature layout template
		layoutPath := filepath.Join(path, "templates", "layout.html")
		ts, err = ts.ParseFiles(layoutPath)
		if err != nil {
			return fmt.Errorf("failed to parse layout template: %w", err)
		}

		// Finally add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return fmt.Errorf("failed to parse page template: %w", err)
		}

		// Store in cache with feature-prefixed name
		cacheKey := fmt.Sprintf("%s_%s", name, strings.TrimSuffix(pageName, ".html"))
		m.pageTemplates[cacheKey] = ts
		log.Printf("  Cached as: %s", cacheKey)
	}

	// Register any partials for HTMX etc.
	partialsPath := filepath.Join(path, "templates", "partials")
	if _, err := os.Stat(partialsPath); os.IsNotExist(err) {
		// Directory doesn't exist
		log.Printf("No partials directory for feature %s", name)
	} else {
		// Directory exists, check for files
		partials, err := filepath.Glob(filepath.Join(partialsPath, "*.html"))
		if err != nil {
			return fmt.Errorf("error checking for partials: %w", err)
		}

		if len(partials) == 0 {
			log.Printf("Partials directory exists but no .html files found for feature %s", name)
		} else {
			// Parse all partials for this feature
			ts, err := template.ParseFiles(partials...)
			if err != nil {
				return fmt.Errorf("failed to parse partial templates: %w", err)
			}
			m.partialTemplates[name] = ts
			log.Printf("  Cached %d partials for feature: %s", len(partials), name)
		}
	}

	// Store nav items
	m.navItems = append(m.navItems, navItems...)

	return nil
}

func (m *Manager) Render(w http.ResponseWriter, r *http.Request, feature, page string, data interface{}) error {
	templateName := fmt.Sprintf("%s_%s", feature, page)
	log.Printf("Looking for template: %s", templateName)

	ts, ok := m.pageTemplates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	viewData := struct {
		Data     interface{}
		NavItems []interfaces.NavItem
		Feature  string
		Error    string
	}{
		Data:     data,
		NavItems: m.navItems,
		Feature:  feature,
		Error:    r.URL.Query().Get("error"),
	}

	// Use buffer for atomic writes
	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "base", viewData); err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}

	buf.WriteTo(w)
	return nil
}

func (m *Manager) RenderPartial(w http.ResponseWriter, r *http.Request, feature string, partial string, data interface{}) error {
	log.Printf("RenderPartial: feature=%s, partial=%s", feature, partial)

	ts, ok := m.partialTemplates[feature]
	if !ok {
		log.Printf("Feature %s not found in partials", feature)
		return fmt.Errorf("feature %s not found", feature)
	}

	viewData := struct {
		Data interface{}
	}{
		Data: data,
	}

	// Use buffer for atomic writes
	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, partial, viewData); err != nil {
		return err
	}

	buf.WriteTo(w)
	return nil
}
