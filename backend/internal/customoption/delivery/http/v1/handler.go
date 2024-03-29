package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	hu "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v "github.com/go-ozzo/ozzo-validation"
)

type CustomOptionHandler struct {
	router http.Handler
	uc     domain.CustomOptionUsecase
}

func NewCustomOptionHandler(uc domain.CustomOptionUsecase) *CustomOptionHandler {
	router := chi.NewRouter()
	handler := &CustomOptionHandler{router: router, uc: uc}

	router.Get("/", handler.GetCustomOptions)
	router.Get("/{id}", handler.GetCustomOptionById)
	router.Post("/", handler.CreateCustomOption)
	router.Put("/{id}", handler.UpdateCustomOption)
	router.Delete("/{id}", handler.DeleteCustomOption)

	return handler
}

func (h *CustomOptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

type customOptionResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *CustomOptionHandler) GetCustomOptions(w http.ResponseWriter, r *http.Request) {
	filter, err := h.getFilter(r.URL.Query())
	if err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("parse filter error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	customOptions, err := h.uc.GetCustomOptions(r.Context(), filter)
	if err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("get custom options error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	customOptionResponses := make([]*customOptionResponse, len(customOptions))
	for i, co := range customOptions {
		customOptionResponses[i] = toCustomOptionResponse(co)
	}

	hu.SuccessResponse(w, r, customOptionResponses)
}

func (h *CustomOptionHandler) GetCustomOptionById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	customOption, err := h.uc.GetCustomOptionById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("get custom option error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, toCustomOptionResponse(customOption))
}

type createCustomOptionInput struct {
	Name string `json:"name"`
}

func (ci *createCustomOptionInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
	)
}

func (h *CustomOptionHandler) CreateCustomOption(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	var input createCustomOptionInput
	if err := render.Bind(r, &input); err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	err := h.uc.CreateCustomOption(r.Context(), &domain.CustomOption{Name: input.Name})
	if err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("create custom option error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

type updateCustomOptionInput struct {
	Name string `json:"name"`
}

func (ci *updateCustomOptionInput) Bind(r *http.Request) error {
	return v.ValidateStruct(ci,
		v.Field(&ci.Name, v.Required, v.Length(1, 50)),
	)
}

func (h *CustomOptionHandler) UpdateCustomOption(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	id := chi.URLParam(r, "id")

	var input updateCustomOptionInput
	if err := render.Bind(r, &input); err != nil {
		hu.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	err := h.uc.UpdateCustomOption(r.Context(), id, &domain.CustomOption{
		Name: input.Name,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("update custom option error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

func (h *CustomOptionHandler) DeleteCustomOption(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteCustomOption(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		hu.FailureResponse(
			w, r,
			fmt.Errorf("delete custom option error - %w", err),
			status,
		)
		return
	}

	hu.SuccessResponse(w, r, nil)
}

func (h *CustomOptionHandler) getFilter(params url.Values) (*domain.CustomOptionFilter, error) {
	var limit int
	var offset int
	var name string

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

	name = params.Get("name")

	return domain.NewCustomOptionFilter(limit, offset, name)
}

func toCustomOptionResponse(customOption *domain.CustomOption) *customOptionResponse {
	return &customOptionResponse{
		Id:   customOption.Id,
		Name: customOption.Name,
	}
}
