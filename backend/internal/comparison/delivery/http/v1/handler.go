package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	hu "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ComparisonUsecase interface {
	Comparisons(ctx context.Context, filter domain.ComparisonFilter) ([]domain.Comparison, error)
	ComparisonById(ctx context.Context, id string) (domain.Comparison, error)
	UpdateComparison(ctx context.Context, id string, comparison domain.Comparison) error
	CreateComparison(ctx context.Context, comparison domain.Comparison) error
	DeleteComparison(ctx context.Context, id string) error
}

type ComparisonHandler struct {
	router http.Handler
	uc     ComparisonUsecase
}

func NewComparisonHandler(uc ComparisonUsecase) *ComparisonHandler {
	router := chi.NewRouter()
	handler := &ComparisonHandler{router: router, uc: uc}

	router.Get("/", handler.Comparisons)
	router.Get("/{id}", handler.ComparisonById)
	router.Post("/", handler.CreateComparison)
	router.Put("/{id}", handler.UpdateComparison)
	router.Delete("/{id}", handler.DeleteComparison)

	return handler
}

func (h *ComparisonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

type comparisonResponse struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	CustomOptionIds []string  `json:"custom_option_ids"`
}

func (h *ComparisonHandler) Comparisons(w http.ResponseWriter, r *http.Request) {
	filter, err := h.getFilter(r.URL.Query())
	if err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("parse filter error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	comparisons, err := h.uc.Comparisons(r.Context(), filter)
	if err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("get comparisons error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	comparisonResponses := make([]comparisonResponse, len(comparisons))
	for i, c := range comparisons {
		comparisonResponses[i] = toComparisonResponse(c)
	}

	hu.SuccessResponse(w, r, comparisonResponses)
}

func (h *ComparisonHandler) ComparisonById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	comparison, err := h.uc.ComparisonById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("get comparison error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, toComparisonResponse(comparison))
}

type createComparisonInput struct {
	Name            string   `json:"name"`
	CustomOptionIds []string `json:"custom_option_ids"`
}

func (ci *createComparisonInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
		v.Field(&ci.CustomOptionIds, v.Each(is.UUIDv4)),
	)
}

func (h *ComparisonHandler) CreateComparison(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	var input createComparisonInput
	if err := render.Bind(r, &input); err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	if input.CustomOptionIds == nil {
		input.CustomOptionIds = make([]string, 0)
	}

	err := h.uc.CreateComparison(r.Context(), domain.Comparison{
		Name:            input.Name,
		CustomOptionIds: input.CustomOptionIds,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrAlreadyExists) {
			status = http.StatusBadRequest
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("create comparison error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

type updateComparisonInput struct {
	Name            string   `json:"name"`
	CustomOptionIds []string `json:"custom_option_ids"`
}

func (ci *updateComparisonInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
		v.Field(&ci.CustomOptionIds, v.Each(is.UUIDv4)),
	)
}

func (h *ComparisonHandler) UpdateComparison(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	id := chi.URLParam(r, "id")

	var input updateComparisonInput
	if err := render.Bind(r, &input); err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	if input.CustomOptionIds == nil {
		input.CustomOptionIds = make([]string, 0)
	}

	err := h.uc.UpdateComparison(r.Context(), id, domain.Comparison{
		Name:            input.Name,
		CustomOptionIds: input.CustomOptionIds,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("update comparison error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

func (h *ComparisonHandler) DeleteComparison(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteComparison(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("delete comparison error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

func (h *ComparisonHandler) getFilter(params url.Values) (domain.ComparisonFilter, error) {
	var limit int
	var offset int
	var orderBy string

	var err error

	limitStr := params.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return domain.ComparisonFilter{}, fmt.Errorf("incorrect limit value")
		}
	}

	offsetStr := params.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return domain.ComparisonFilter{}, fmt.Errorf("incorrect offset value")
		}
	}

	orderBy = params.Get("order_by")

	return domain.NewComparisonFilter(limit, offset, orderBy)
}

func toComparisonResponse(comparison domain.Comparison) comparisonResponse {
	return comparisonResponse{
		Id:              comparison.Id,
		Name:            comparison.Name,
		CreatedAt:       comparison.CreatedAt,
		CustomOptionIds: comparison.CustomOptionIds,
	}
}
