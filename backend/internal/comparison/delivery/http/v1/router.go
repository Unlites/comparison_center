package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	httputils "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
)

type ComparisonHandler struct {
	Router http.Handler
	uc     domain.ComparisonUsecase
}

func NewComparisonHandler(uc domain.ComparisonUsecase) *ComparisonHandler {
	router := chi.NewRouter()
	handler := &ComparisonHandler{Router: router, uc: uc}

	router.Get("/", handler.getComparisons)
	router.Get("/{id}", handler.getComparisonById)
	router.Post("/", handler.createComparison)
	router.Put("/{id}", handler.updateComparison)
	router.Delete("/{id}", handler.deleteComparison)

	return handler
}

func (h *ComparisonHandler) getComparisons(w http.ResponseWriter, r *http.Request) {
	filter, err := h.getFilter(r.URL.Query())
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("parse filter error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	comparisons, err := h.uc.GetComparisons(r.Context(), filter)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("get comparisons error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	httputils.SuccessResponse(w, r, comparisons)
}

func (h *ComparisonHandler) getComparisonById(w http.ResponseWriter, r *http.Request) {

}

func (h *ComparisonHandler) updateComparison(w http.ResponseWriter, r *http.Request) {

}

func (h *ComparisonHandler) createComparison(w http.ResponseWriter, r *http.Request) {

}

func (h *ComparisonHandler) deleteComparison(w http.ResponseWriter, r *http.Request) {

}

func (h *ComparisonHandler) getFilter(params url.Values) (*domain.ComparisonFilter, error) {
	var limit int
	var offset int
	var orderBy string

	var err error

	limitStr := params.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect limit value")
		}
	}

	offsetStr := params.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("incorrect offset value")
		}
	}

	orderBy = params.Get("order_by")

	return domain.NewComparisonFilter(limit, offset, orderBy)
}
