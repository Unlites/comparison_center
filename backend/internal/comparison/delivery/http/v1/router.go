package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	httputils "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v "github.com/go-ozzo/ozzo-validation"
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
	id := chi.URLParam(r, "id")

	comparison, err := h.uc.GetComparisonById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("get comparisons error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, comparison)
}

type createComparisonInput struct {
	Name string `json:"name"`
}

func (ci *createComparisonInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
	)
}

func (h *ComparisonHandler) createComparison(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	var input createComparisonInput
	if err := render.Bind(r, &input); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	err := h.uc.CreateComparison(r.Context(), &domain.Comparison{Name: input.Name})
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("create comparison error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
}

type updateComparisonInput struct {
	Name            string   `json:"name"`
	CustomOptionIds []string `json:"custom_option_ids"`
}

func (ci *updateComparisonInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
		v.Field(&ci.CustomOptionIds, v.Required),
	)
}

func (h *ComparisonHandler) updateComparison(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	id := chi.URLParam(r, "id")

	var input updateComparisonInput
	if err := render.Bind(r, &input); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	err := h.uc.UpdateComparison(r.Context(), id, &domain.Comparison{
		Name:            input.Name,
		CustomOptionIds: input.CustomOptionIds,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("update comparison error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
}

func (h *ComparisonHandler) deleteComparison(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteComparison(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("update comparison error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
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
