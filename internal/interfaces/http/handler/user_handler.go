package handler

import (
	"encoding/json"
	"net/http"
	"trainer/internal/application/dto"
	"trainer/internal/application/usecase"
	"trainer/internal/interfaces/http/response"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	createUserUC *usecase.CreateUser
	updateUserUC *usecase.UpdateUser
	deleteUserUC *usecase.DeleteUser
	getUserUC    *usecase.GetUser
	listUserUC   *usecase.ListUser
}

func NewUserHandler(
	createUserUC *usecase.CreateUser,
	updateUserUC *usecase.UpdateUser,
	deleteUserUC *usecase.DeleteUser,
	getUserUC *usecase.GetUser,
	listUserUC *usecase.ListUser,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		updateUserUC: updateUserUC,
		deleteUserUC: deleteUserUC,
		getUserUC:    getUserUC,
		listUserUC:   listUserUC,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	userResp, err := h.createUserUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, userResp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	vars := mux.Vars(r)
	req.Id = vars["id"]

	userResp, err := h.updateUserUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, userResp)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	vars := mux.Vars(r)
	req.Id = vars["id"]

	err := h.deleteUserUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, struct{}{})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var req dto.GetUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	vars := mux.Vars(r)
	req.Id = vars["id"]

	userResp, err := h.getUserUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, userResp)
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	var req dto.ListUserRequest

	userResp, err := h.listUserUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, userResp)
}
