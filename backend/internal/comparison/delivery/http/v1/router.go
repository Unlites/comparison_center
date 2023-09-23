package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ComparisonHandler struct {
	Router http.Handler
	// uc domain.ComparisonUsecase
}

func NewComparisonHandler() *ComparisonHandler {
	router := chi.NewRouter()

	router.Get("/ping", ping)

	return &ComparisonHandler{Router: router}
}

func ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]interface{}{
		"message": "pong from comparison handler",
	})
}
