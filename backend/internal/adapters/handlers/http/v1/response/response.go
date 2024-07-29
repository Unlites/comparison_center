package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, data any) {
	render.JSON(w, r, &response{
		Success: true,
		Data:    data,
	})
}

func FailureResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	render.Status(r, statusCode)
	render.JSON(w, r, &response{
		Success: false,
		Message: err.Error(),
	})
}
