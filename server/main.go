package server

import (
	"embed"
	"html/template"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mehix/vf-sec-ctrls/pkg/domain/categ/memory"
	"github.com/mehix/vf-sec-ctrls/pkg/service/categ"
)

//go:embed templates
var htmlTmpl embed.FS

type Server struct {
	tmpl         *template.Template
	tokenAuth    *jwtauth.JWTAuth
	categService *categ.Service
}

// New returns a new Server with an in-memory categories database
func New() *Server {
	t, err := template.ParseFS(htmlTmpl, "**/*.tmpl")
	if err != nil {
		panic(err)
	}

	signingKey := "lakjsdk%(asjfkajFKJd4%&safkjadkfju98798ID*94823%&)f87"
	if s, ok := os.LookupEnv("JWT_SECRET"); ok && s != "" {
		signingKey = s
	}

	return &Server{
		tmpl:         t,
		tokenAuth:    jwtauth.New("HS256", []byte(signingKey), nil),
		categService: categ.NewService(memory.NewRepository()),
	}
}
