package chi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type defaultRouter struct {
	Handler chi.Router
}

func NewDefaultRouter() *defaultRouter {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RedirectSlashes)

	return &defaultRouter{Handler: router}
}

func (dr *defaultRouter) RegisterHandlers(version string, handlers map[string]http.Handler) {
	versionPrefix := fmt.Sprintf("/api/%s/", version)
	for prefix, handler := range handlers {
		dr.Handler.Mount(versionPrefix+prefix, handler)
	}
}
