package users

import (
	"net/http"
	"time"

	"github.com/MickDuprez/gobase/internal/core/auth"
	"github.com/MickDuprez/gobase/internal/core/interfaces"
)

type Handler struct {
	app interfaces.App
}

func (h *Handler) loginForm(w http.ResponseWriter, r *http.Request) {
	h.app.RenderTemplate(w, r, "users", "login", nil)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.app.Auth().ValidateUser(email, password)
	if err != nil {
		// Redirect back to login with error
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	// Create session
	session, err := h.app.Auth().CreateSession(user.ID, 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) registerForm(w http.ResponseWriter, r *http.Request) {
	h.app.RenderTemplate(w, r, "users", "register", nil)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	user, err := h.app.Auth().CreateUser(email, password, name)
	if err != nil {
		// Handle registration errors (e.g., duplicate email)
		http.Redirect(w, r, "/register?error=registration_failed", http.StatusSeeOther)
		return
	}

	// Auto-login after registration
	session, err := h.app.Auth().CreateSession(user.ID, 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete session from database
		h.app.Auth().DeleteSession(cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (added by auth middleware)
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		User *auth.User
	}{
		User: user,
	}

	h.app.RenderTemplate(w, r, "users", "profile", data)
}

func (h *Handler) addProfileInfo(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "user not found!", http.StatusInternalServerError)
		return
	}

	h.app.RenderPartial(w, r, "users", "profile_form", nil)
}

func (h *Handler) showProfileInfo(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "user not found!", http.StatusInternalServerError)
		return
	}

	// get data from database here

	// set the data (hard coded for now)
	info := map[string]string{
		"location": "my place",
		"bio":      "I do stuff",
		"website":  "www.micko.com",
	}

	h.app.RenderPartial(w, r, "users", "profile_info", info)

}

func (h *Handler) saveProfileInfo(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "user not found!", http.StatusInternalServerError)
		return
	}

	// parse the form
	r.ParseForm()
	info := map[string]string{
		"location": r.FormValue("location"),
		"bio":      r.FormValue("bio"),
		"website":  r.FormValue("website"),
	}

	// save to db here...

	h.app.RenderPartial(w, r, "users", "profile_info", info)
}
