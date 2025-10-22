package response

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, ErrorResponse{Error: err.Error()})
}

func BadRequest(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, err)
}

func NotFound(w http.ResponseWriter, err error) {
	Error(w, http.StatusNotFound, err)
}

func InternalError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, err)
}
