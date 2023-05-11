package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/mehix/vf-sec-ctrls/categ"
)

var h *categ.Hierarchy

func (s *Server) Handlers(orig *categ.Hierarchy) http.Handler {
	h = orig
	m := chi.NewMux()

	m.Use(middleware.Logger)
	m.Use(middleware.RequestID)
	m.Use(middleware.RealIP)
	m.Use(middleware.Maybe(middleware.URLFormat, func(r *http.Request) bool { return !strings.HasPrefix(r.URL.Path, "/controls/") }))

	m.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	m.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html;charset=utf-8")
		if err := s.tmpl.ExecuteTemplate(w, "login", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	m.Group(func(m chi.Router) {
		m.Use(jwtauth.Verifier(s.tokenAuth))
		m.Use(jwtauth.Authenticator)

		m.Get("/user", s.handleShowUser)
		m.Get("/upload", s.handleShowUploadForm)
		m.Post("/upload", s.handleUpload)
	})

	m.Post("/login", s.handleLogin)

	m.Get("/controls", func(w http.ResponseWriter, r *http.Request) {
		format, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)

		entries := make([]categ.Entry, 0)
		if h != nil {
			entries = h.Entries()
		}

		switch format {
		case "json":
			renderJson(w, r, entries)
		case "txt":
			renderTxt(w, r, entries)
		default:
			renderJson(w, r, entries)
		}

	})
	m.Get("/controls/{format:(json|txt)}/{categoryID:[0-9.]+}", func(w http.ResponseWriter, r *http.Request) {
		categoryID := chi.URLParam(r, "categoryID")

		controls := make([]categ.Entry, 0)
		if h != nil {
			controls = h.ControlsByCategory(categoryID)
		}

		format := chi.URLParam(r, "format")

		fmt.Printf("Requested format: %s\n", format)

		switch format {
		case "json":
			renderJson(w, r, controls)
		case "txt":
			renderTxt(w, r, controls)
		default:
			renderJson(w, r, controls)
		}

	})

	return m

}

func renderJson(w http.ResponseWriter, r *http.Request, entries []categ.Entry) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTxt(w http.ResponseWriter, r *http.Request, entries []categ.Entry) {
	w.Header().Set("Context-type", "text/plain;charset=utf-8")
	for _, e := range entries {
		w.Write([]byte(e.String() + "\n"))
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	passwd := r.FormValue("passwd")

	// TODO: do the actual validation
	if user != "mihai" || passwd != "mihaivf" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := map[string]any{
		"user_id":    user,
		"issuer":     "vattenfall",
		"session_id": uuid.NewString(),
	}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, tokenStr, err := s.tokenAuth.Encode(claims)
	if err != nil {
		log.Printf("creating jwt: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:  "jwt",
		Value: tokenStr,
		Path:  "/",
		//SameSite: http.SameSiteStrictMode,
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	userAgent := r.Header.Get("User-Agent")
	if strings.HasPrefix(userAgent, "curl") {
		// return token to CURL
		w.Write([]byte(tokenStr))
	} else {
		// redirect browser to a secure page
		w.Header().Set("Location", "/upload")
		w.WriteHeader(http.StatusFound)
	}
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uploaded := r.MultipartForm.File
	f, ok := uploaded["file"]
	if !ok {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}

	in, err := f[0].Open()
	if err != nil {
		log.Printf("reading uploaded file: %v\n", err)
		http.Error(w, "cannot open uploaded file", http.StatusInternalServerError)
		return
	}
	defer in.Close()

	wsValue := r.MultipartForm.Value["worksheet"]
	ws := ""
	if len(wsValue) > 0 {
		ws = wsValue[0]
	}
	newH, err := categ.NewFromReader(in, ws)
	if err != nil {
		log.Printf("creating hierarchy: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newH.Print(true)

	h = newH

	w.Header().Set("Location", "/controls.json")
	w.WriteHeader(http.StatusFound)

}

func (s *Server) handleShowUser(w http.ResponseWriter, r *http.Request) {

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Printf("fetching claims from request: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := claims["user_id"].(string)
	sessionID := claims["session_id"].(string)

	w.Header().Set("Content-type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		UserID    string `json:"user_id"`
		SessionID string `json:"session_id"`
	}{
		UserID:    userID,
		SessionID: sessionID,
	})
}

func (s *Server) handleShowUploadForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "upload", nil); err != nil {
		log.Printf("executing template 'upload': %v\n", err)
		http.Error(w, "page is broken", http.StatusInternalServerError)
		return
	}
}
