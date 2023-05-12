package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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

	m.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	m.Get("/login", s.handleShowLogin)
	m.Get("/logout", s.handleLogout)

	m.Post("/login", s.handleLogin)

	m.Group(func(m chi.Router) {
		m.Use(jwtauth.Verifier(s.tokenAuth))
		m.Use(Authenticator)

		m.Get("/user", s.handleShowUser)
		m.Get("/upload", s.handleShowUploadForm)
		m.Post("/upload", s.handleUpload)

		m.Route("/controls", func(m chi.Router) {
			m.Get("/", s.handleShowControlsByID)
			m.Get("/{id:[0-9.]+}", s.handleShowControlsByID)
			m.Get("/{id:[0-9.]+}/{format:(json|txt|html)}", s.handleShowControlsByID)
			m.Get("/list/{format:(json|txt)}", s.handleShowAllControls)
		})
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

func (s *Server) handleShowLogin(w http.ResponseWriter, r *http.Request) {
	respData := struct {
		Error string
	}{
		Error: r.URL.Query().Get("error"),
	}

	w.Header().Set("Content-type", "text/html;charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "login", respData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	passwd := r.FormValue("passwd")

	// TODO: do the actual validation
	if user != "mihai" || passwd != "mihaivf" {
		redirectTo := "/login?error=" + url.QueryEscape("access denied")
		w.Header().Set("Location", redirectTo)
		w.WriteHeader(http.StatusFound)
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

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("jwt")
	if cookie != nil {
		cookie.Expires = time.Now()
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
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

	h = newH

	w.Header().Set("Location", "/controls/")
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

func (s *Server) handleShowControlsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	controls := make([]categ.Entry, 0)
	if h != nil {
		controls = h.ControlsByCategory(id)
	}

	format := chi.URLParam(r, "format")

	renderHTML := func() {
		type responseData struct {
			Controls []categ.Entry
			Control  *categ.Entry // search result
			Parent   *categ.Entry
			Children []categ.Entry
		}

		data := responseData{}
		if h != nil {
			data.Controls = h.Entries()
			data.Control = h.Entry(id)
			data.Parent = h.Parent(id)
			data.Children = h.ControlsByCategory(id)
		}

		w.Header().Set("Content-type", "text/html; charset=utf-8")
		if err := s.tmpl.ExecuteTemplate(w, "controls", data); err != nil {
			log.Printf("executing template 'controls': %v\n", err)
			http.Error(w, "page is broken", http.StatusInternalServerError)
			return
		}
	}

	switch format {
	case "json":
		renderJson(w, r, controls)
	case "txt":
		renderTxt(w, r, controls)
	default:
		renderHTML()
	}

}

func (s *Server) handleShowAllControls(w http.ResponseWriter, r *http.Request) {
	format := chi.URLParam(r, "format")

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

}
