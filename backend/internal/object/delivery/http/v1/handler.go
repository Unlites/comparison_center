package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	httputils "github.com/Unlites/comparison_center/backend/internal/utils/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ObjectHandler struct {
	Router        http.Handler
	MaxUploadSize int64
	PhotosDir     string
	uc            domain.ObjectUsecase
}

func NewObjectHandler(uc domain.ObjectUsecase, photosDir string, maxSize int64) *ObjectHandler {
	router := chi.NewRouter()
	handler := &ObjectHandler{
		Router:        router,
		MaxUploadSize: maxSize << 20,
		PhotosDir:     photosDir,
		uc:            uc,
	}

	router.Get("/", handler.getObjects)
	router.Get("/{id}", handler.getObjectById)
	router.Post("/", handler.createObject)
	router.Put("/{id}", handler.updateObject)
	router.Delete("/{id}", handler.deleteObject)

	router.Get("/{id}/photo", handler.getObjectPhoto)
	router.Post("/{id}/photo", handler.uploadObjectPhoto)

	return handler
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

func (h *ObjectHandler) getObjects(w http.ResponseWriter, r *http.Request) {
	filter, err := h.getFilter(r.URL.Query())
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("parse filter error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	objects, err := h.uc.GetObjects(r.Context(), filter)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("get objects error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	objectResponses := make([]*objectResponse, len(objects))
	for i, o := range objects {
		objectResponses[i] = toObjectResponse(o)
	}

	httputils.SuccessResponse(w, r, objectResponses)
}

func (h *ObjectHandler) getObjectById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	object, err := h.uc.GetObjectById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("get object error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, toObjectResponse(object))
}

type createObjectInput struct {
	Name         string `json:"name"`
	Rating       int    `json:"rating"`
	Advs         string `json:"advs"`
	Disadvs      string `json:"disadvs"`
	ComparisonId string `json:"comparison_id"`
}

func (oi *createObjectInput) Bind(r *http.Request) error {
	return v.ValidateStruct(oi,
		v.Field(&oi.Name, v.Required, v.Length(1, 50)),
		v.Field(&oi.Rating, v.Min(1), v.Max(10)),
		v.Field(&oi.Advs, v.Length(1, 3000)),
		v.Field(&oi.Disadvs, v.Length(1, 3000)),
		v.Field(&oi.ComparisonId, v.Required, is.UUIDv4),
	)
}

func (h *ObjectHandler) createObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	var input createObjectInput
	if err := render.Bind(r, &input); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	err := h.uc.CreateObject(r.Context(), &domain.Object{
		Name:         input.Name,
		Rating:       input.Rating,
		Advs:         input.Advs,
		Disadvs:      input.Disadvs,
		ComparisonId: input.ComparisonId,
	})
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("create object error - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
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
		v.Field(&oi.Rating, v.Min(1), v.Max(10)),
		v.Field(&oi.Advs, v.Length(1, 3000)),
		v.Field(&oi.Disadvs, v.Length(1, 3000)),
		v.Field(&oi.CustomOptions, v.Each(v.Map(
			v.Key("id", v.Required, is.UUIDv4),
			v.Key("value", v.Required, v.Length(1, 100)),
		))),
	)
}

func (h *ObjectHandler) updateObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - request body required"),
			http.StatusBadRequest,
		)
		return
	}

	id := chi.URLParam(r, "id")

	var input updateObjectInput
	if err := render.Bind(r, &input); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("validation error - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	objCustOpts := make([]*domain.ObjectCustomOption, len(input.CustomOptions))
	for i, opt := range input.CustomOptions {
		objCustOpts[i] = &domain.ObjectCustomOption{
			ObjectId:       id,
			CustomOptionId: opt["id"],
			Value:          opt["value"],
		}
	}

	err := h.uc.UpdateObject(r.Context(), id, &domain.Object{
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

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("update object error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
}

func (h *ObjectHandler) deleteObject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.uc.DeleteObject(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("delete object error - %w", err),
			status,
		)
		return
	}

	httputils.SuccessResponse(w, r, nil)
}

func (h *ObjectHandler) uploadObjectPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, h.MaxUploadSize)
	if err := r.ParseMultipartForm(h.MaxUploadSize); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("too big size, max is %d", h.MaxUploadSize),
			http.StatusBadRequest,
		)
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		httputils.FailureResponse(
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
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to read photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("invalid photo format, must be image/jpeg or image/jpeg"),
			http.StatusBadRequest,
		)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to seek photo - %w", err),
			http.StatusBadRequest,
		)
		return
	}

	photoPath := fmt.Sprintf(
		"%s/%s%s",
		h.PhotosDir,
		uuid.NewString(),
		filepath.Ext(fileHeader.Filename),
	)

	newPhoto, err := os.Create(photoPath)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to create photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}
	defer newPhoto.Close()

	_, err = io.Copy(newPhoto, file)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to save photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	if err := h.uc.SetObjectPhotoPath(r.Context(), id, photoPath); err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to set photo path - %w", err),
			http.StatusInternalServerError,
		)
		os.Remove(photoPath)
		return
	}

	httputils.SuccessResponse(w, r, nil)
}
func (h *ObjectHandler) getObjectPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	object, err := h.uc.GetObjectById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}

		httputils.FailureResponse(
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

		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to open photo - %w", err),
			status,
		)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		httputils.FailureResponse(
			w, r,
			fmt.Errorf("failed to send photo - %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
}

func (h *ObjectHandler) getFilter(params url.Values) (*domain.ObjectFilter, error) {
	var limit int
	var offset int
	var name string
	var comparisonId string

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

	orderBy := params.Get("order_by")
	name = params.Get("name")
	comparisonId = params.Get("comparison_id")

	return domain.NewObjectFilter(limit, offset, orderBy, name, comparisonId)
}

func toObjectResponse(object *domain.Object) *objectResponse {
	customOpts := make([]map[string]string, len(object.ObjectCustomOptions))
	for i, co := range object.ObjectCustomOptions {
		customOpts[i] = map[string]string{
			"id":    co.CustomOptionId,
			"value": co.Value,
		}
	}
	return &objectResponse{
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
