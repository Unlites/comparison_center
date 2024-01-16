package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	httputils "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ComparisonHandler struct {
	router http.Handler
	uc     domain.ComparisonUsecase
}

func NewComparisonHandler(uc domain.ComparisonUsecase) *ComparisonHandler {
	router := chi.NewRouter()
	handler := &ComparisonHandler{router: router, uc: uc}

	router.Get("/", handler.GetComparisons)
	router.Get("/{id}", handler.GetComparisonById)
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

func (h *ComparisonHandler) GetComparisons(w http.ResponseWriter, r *http.Request) {
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

	comparisonResponses := make([]*comparisonResponse, len(comparisons))
	for i, c := range comparisons {
		comparisonResponses[i] = toComparisonResponse(c)
	}

	httputils.SuccessResponse(w, r, comparisonResponses)
}

func (h *ComparisonHandler) GetComparisonById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	comparison, err := h.uc.GetComparisonById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("get comparison error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, toComparisonResponse(comparison))
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

	if input.CustomOptionIds == nil {
		input.CustomOptionIds = make([]string, 0)
	}

	err := h.uc.CreateComparison(r.Context(), &domain.Comparison{
		Name:            input.Name,
		CustomOptionIds: input.CustomOptionIds,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrAlreadyExists) {
			status = http.StatusBadRequest
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("create comparison error - %w", err),
			status,
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
		v.Field(&ci.CustomOptionIds, v.Each(is.UUIDv4)),
	)
}

func (h *ComparisonHandler) UpdateComparison(w http.ResponseWriter, r *http.Request) {
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

	if input.CustomOptionIds == nil {
		input.CustomOptionIds = make([]string, 0)
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

func (h *ComparisonHandler) DeleteComparison(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteComparison(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("delete comparison error - %w", err),
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

func toComparisonResponse(comparison *domain.Comparison) *comparisonResponse {
	return &comparisonResponse{
		Id:              comparison.Id,
		Name:            comparison.Name,
		CreatedAt:       comparison.CreatedAt,
		CustomOptionIds: comparison.CustomOptionIds,
	}
}
