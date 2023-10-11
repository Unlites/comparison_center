package chi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Handler chi.Router
}

func NewDefaultRouter() *Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RedirectSlashes)

	return &Router{Handler: router}
}

func (r *Router) RegisterHandlers(version string, handlers map[string]http.Handler) {
	versionPrefix := fmt.Sprintf("/api/%s/", version)
	for prefix, handler := range handlers {
		r.Handler.Mount(versionPrefix+prefix, handler)
	}
}
