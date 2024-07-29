package object

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/response"
	"github.com/Unlites/comparison_center/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ObjectUsecase interface {
	GetObjects(ctx context.Context, filter domain.ObjectFilter) ([]domain.Object, error)
	GetObjectById(ctx context.Context, id string) (domain.Object, error)
	UpdateObject(ctx context.Context, id string, object domain.Object) error
	CreateObject(ctx context.Context, object domain.Object) (string, error)
	DeleteObject(ctx context.Context, id string) error
	SetObjectPhotoPath(ctx context.Context, id, path string) error
}

type ObjectHandler struct {
	router        http.Handler
	maxUploadSize int64
	photosDir     string
	uc            ObjectUsecase
}

func NewObjectHandler(uc ObjectUsecase, photosDir string, maxSize int64) *ObjectHandler {
	router := chi.NewRouter()
	handler := &ObjectHandler{
		router:        router,
		maxUploadSize: maxSize << 20,
		photosDir:     photosDir,
		uc:            uc,
	}

	router.Get("/", handler.GetObjects)
	router.Get("/{id}", handler.GetObjectById)
	router.Post("/", handler.CreateObject)
	router.Put("/{id}", handler.UpdateObject)
	router.Delete("/{id}", handler.DeleteObject)

	router.Get("/{id}/photo", handler.GetObjectPhoto)
	router.Post("/{id}/photo", handler.UploadObjectPhoto)

	return handler
}

func (h *ObjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

type objectResponse struct {
	Id            string              `json:"id"`
	Name          string              `json:"name"`
	Rating        int                 `json:"rating"`
	CreatedAt     time.Time           `json:"created_at"`
	Advs          string              `json:"advs"`
	Disadvs       string              `json:"disadvs"`
	ComparisonId  string              `json:"comparison_id"`
	CustomOptions []map[string]string `json:"custom_options"`
}

func (h *ObjectHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	filter, err := h.getFilter(r.URL.Query())
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("parse filter error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	objects, err := h.uc.GetObjects(r.Context(), filter)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("get objects error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	objectResponses := make([]objectResponse, len(objects))
	for i, o := range objects {
		objectResponses[i] = toObjectResponse(o)
	}

	response.SuccessResponse(w, r, objectResponses)
}

func (h *ObjectHandler) GetObjectById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	object, err := h.uc.GetObjectById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		response.FailureResponse(
			w, r,
			fmt.Errorf("get object error - %w", err),
			status,
		)
		return
	}

	response.SuccessResponse(w, r, toObjectResponse(object))
}

type createObjectInput struct {
	Name          string              `json:"name"`
	Rating        int                 `json:"rating"`
	Advs          string              `json:"advs"`
	Disadvs       string              `json:"disadvs"`
	ComparisonId  string              `json:"comparison_id"`
	CustomOptions []map[string]string `json:"custom_options"`
}

func (oi *createObjectInput) Bind(r *http.Request) error {
	return v.ValidateStruct(oi,
		v.Field(&oi.Name, v.Required, v.Length(1, 50)),
		v.Field(&oi.Rating, v.Required, v.Min(1), v.Max(10)),
		v.Field(&oi.Advs, v.Length(1, 3000)),
		v.Field(&oi.Disadvs, v.Length(1, 3000)),
		v.Field(&oi.ComparisonId, v.Required, is.UUIDv4),
		v.Field(&oi.CustomOptions, v.Each(v.Map(
			v.Key("id", v.Required, is.UUIDv4),
			v.Key("value", v.Required, v.Length(1, 100)),
		))),
	)
}

type returnedIdResponse struct {
	Id string `json:"id"`
}

func (h *ObjectHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		response.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	var input createObjectInput
	if err := render.Bind(r, &input); err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	objCustOpts := make([]domain.ObjectCustomOption, len(input.CustomOptions))
	for i, opt := range input.CustomOptions {
		objCustOpts[i] = domain.ObjectCustomOption{
			CustomOptionId: opt["id"],
			Value:          opt["value"],
		}
	}

	id, err := h.uc.CreateObject(r.Context(), domain.Object{
		Name:                input.Name,
		Rating:              input.Rating,
		Advs:                input.Advs,
		Disadvs:             input.Disadvs,
		ComparisonId:        input.ComparisonId,
		ObjectCustomOptions: objCustOpts,
	})
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("create object error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	response.SuccessResponse(w, r, &returnedIdResponse{Id: id})
}

type updateObjectInput struct {
	Name          string              `json:"name"`
	Rating        int                 `json:"rating"`
	Advs          string              `json:"advs"`
	Disadvs       string              `json:"disadvs"`
	CustomOptions []map[string]string `json:"custom_options"`
}

func (oi *updateObjectInput) Bind(r *http.Request) error {
	return v.ValidateStruct(oi,
		v.Field(&oi.Name, v.Required, v.Length(1, 50)),
		v.Field(&oi.Rating, v.Required, v.Min(1), v.Max(10)),
		v.Field(&oi.Advs, v.Length(1, 3000)),
		v.Field(&oi.Disadvs, v.Length(1, 3000)),
		v.Field(&oi.CustomOptions, v.Each(v.Map(
			v.Key("id", v.Required, is.UUIDv4),
			v.Key("value", v.Required, v.Length(1, 100)),
		))),
	)
}

func (h *ObjectHandler) UpdateObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		response.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	id := chi.URLParam(r, "id")

	var input updateObjectInput
	if err := render.Bind(r, &input); err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	objCustOpts := make([]domain.ObjectCustomOption, len(input.CustomOptions))
	for i, opt := range input.CustomOptions {
		objCustOpts[i] = domain.ObjectCustomOption{
			ObjectId:       id,
			CustomOptionId: opt["id"],
			Value:          opt["value"],
		}
	}

	err := h.uc.UpdateObject(r.Context(), id, domain.Object{
		Name:                input.Name,
		Rating:              input.Rating,
		Advs:                input.Advs,
		Disadvs:             input.Disadvs,
		ObjectCustomOptions: objCustOpts,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		response.FailureResponse(
			w, r,
			fmt.Errorf("update object error - %w", err),
			status,
		)
		return
	}

	response.SuccessResponse(w, r, nil)
}

func (h *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteObject(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		response.FailureResponse(
			w, r,
			fmt.Errorf("delete object error - %w", err),
			status,
		)
		return
	}

	response.SuccessResponse(w, r, nil)
}

func (h *ObjectHandler) UploadObjectPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadSize)
	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failted to parse multipart form - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to get photo - %w", err),
			http.StatusBadRequest,
		)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to read photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		response.FailureResponse(
			w, r,
			fmt.Errorf("invalid photo format, must be image/jpeg or image/png"),
			http.StatusBadRequest,
		)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to seek photo - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	photoPath := fmt.Sprintf(
		"%s/%s%s",
		h.photosDir,
		uuid.NewString(),
		filepath.Ext(fileHeader.Filename),
	)

	newPhoto, err := os.Create(photoPath)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to create photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}
	defer newPhoto.Close()

	_, err = io.Copy(newPhoto, file)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to save photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	if err := h.uc.SetObjectPhotoPath(r.Context(), id, photoPath); err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to set photo path - %w", err),
			http.StatusInternalServerError,
		)
		os.Remove(photoPath)
		return
	}

	response.SuccessResponse(w, r, nil)
}
func (h *ObjectHandler) GetObjectPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	object, err := h.uc.GetObjectById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to set photo path - %w", err),
			status,
		)
		return
	}

	file, err := os.Open(object.PhotoPath)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}

		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to open photo - %w", err),
			status,
		)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		response.FailureResponse(
			w, r,
			fmt.Errorf("failed to send photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
}

func (h *ObjectHandler) getFilter(params url.Values) (domain.ObjectFilter, error) {
	var limit int
	var offset int
	var name string
	var comparisonId string

	var err error

	limitStr := params.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return domain.ObjectFilter{}, fmt.Errorf("incorrect limit value")
		}
	}

	offsetStr := params.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return domain.ObjectFilter{}, fmt.Errorf("incorrect offset value")
		}
	}

	orderBy := params.Get("order_by")
	name = params.Get("name")
	comparisonId = params.Get("comparison_id")

	return domain.NewObjectFilter(limit, offset, orderBy, name, comparisonId)
}

func toObjectResponse(object domain.Object) objectResponse {
	customOpts := make([]map[string]string, len(object.ObjectCustomOptions))
	for i, co := range object.ObjectCustomOptions {
		customOpts[i] = map[string]string{
			"id":    co.CustomOptionId,
			"value": co.Value,
		}
	}
	return objectResponse{
		Id:            object.Id,
		Name:          object.Name,
		Rating:        object.Rating,
		CreatedAt:     object.CreatedAt,
		Advs:          object.Advs,
		Disadvs:       object.Disadvs,
		ComparisonId:  object.ComparisonId,
		CustomOptions: customOpts,
	}
}
