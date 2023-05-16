package version

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	Version   = "dev"
	BuildDate = "Somewhere in time"
)

func Append(r *chi.Mux) {
	r.Get("/api/version", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Version   string `json:"version"`
			BuildDate string `json:"build_time"`
		}{
			Version:   Version,
			BuildDate: BuildDate,
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
