package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MickDuprez/gobase/internal/core/utils"
)

type SecurityConfig struct {
	ScriptSources  []string
	StyleSources   []string
	FontSources    []string
	ImageSources   []string
	ConnectSources []string
	DefaultHeaders map[string]string
	IsDevelopment  bool
}

func NewDevSecurityConfig() *SecurityConfig {
	isDev := utils.GetEnvBool("IS_DEV", true)

	cfg := &SecurityConfig{
		ScriptSources: []string{
			"'self'",
			"'unsafe-inline'",
			"https://unpkg.com",
			"https://cdn.jsdelivr.net",
		},
		StyleSources: []string{
			"'self'",
			"'unsafe-inline'",
			"https://cdn.jsdelivr.net",
		},
		FontSources: []string{
			"'self'",
			"https://cdn.jsdelivr.net",
		},
		ImageSources: []string{
			"'self'",
			"data:",
			"https:",
		},
		ConnectSources: []string{
			"'self'",
		},
		DefaultHeaders: map[string]string{
			"Server":                 "",
			"X-Powered-By":           "",
			"X-Content-Type-Options": "nosniff",
			"X-Frame-Options":        "DENY",
			"X-XSS-Protection":       "1; mode=block",
			"Referrer-Policy":        "strict-origin-when-cross-origin",
		},
		IsDevelopment: isDev,
	}

	// Add development-specific settings
	if isDev {
		// Allow WebSocket connections in dev
		if utils.GetEnvBool("ALLOW_WEBSOCKETS", true) {
			cfg.ConnectSources = append(cfg.ConnectSources, "ws:", "wss:")
		}
	} else {
		// Add production-only headers
		cfg.DefaultHeaders["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	}

	return cfg
}

func NewProdSecurityConfig() *SecurityConfig {
	config := &SecurityConfig{
		ScriptSources: []string{
			"'self'",
			"'unsafe-inline'",
			"https://unpkg.com",
			"https://cdn.jsdelivr.net",
		},
		StyleSources: []string{
			"'self'",
			"'unsafe-inline'",
			"https://cdn.jsdelivr.net",
		},
		FontSources: []string{
			"'self'",
			"https://cdn.jsdelivr.net",
		},
		ImageSources: []string{
			"'self'",
			"data:",
			"https:",
		},
		ConnectSources: []string{
			"'self'",
		},
		DefaultHeaders: map[string]string{
			"Server":                    "",
			"X-Powered-By":              "",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "DENY",
			"X-XSS-Protection":          "1; mode=block",
			"Referrer-Policy":           "strict-origin-when-cross-origin",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains", // Only in prod
		},
		IsDevelopment: false,
	}
	return config
}

func (c *SecurityConfig) BuildCSP() string {
	csp := []string{
		"default-src 'self'",
		fmt.Sprintf("script-src %s", strings.Join(c.ScriptSources, " ")),
		fmt.Sprintf("style-src %s", strings.Join(c.StyleSources, " ")),
		fmt.Sprintf("font-src %s", strings.Join(c.FontSources, " ")),
		fmt.Sprintf("img-src %s", strings.Join(c.ImageSources, " ")),
		fmt.Sprintf("connect-src %s", strings.Join(c.ConnectSources, " ")),
	}
	return strings.Join(csp, "; ")
}

func SecurityHeaders(config *SecurityConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Set default headers
			for key, value := range config.DefaultHeaders {
				w.Header().Set(key, value)
			}

			// Set CSP
			w.Header().Set("Content-Security-Policy", config.BuildCSP())

			next.ServeHTTP(w, r)
		}
	}
}

// Time-based brute force protection
func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	// Implement rate limiting logic
	return nil
}

// CSRF protection
func CSRF(next http.HandlerFunc) http.HandlerFunc {
	// Implement CSRF protection
	return nil
}
