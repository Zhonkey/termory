package handler

import (
	"encoding/json"
	"net/http"
	"trainer/internal/application/dto"
	"trainer/internal/application/usecase"
	"trainer/internal/interfaces/http/response"
)

type AuthHandler struct {
	accessTokenUC  *usecase.AccessToken
	refreshTokenUC *usecase.RefreshToken
}

func NewAuthTokenHandler(loginUserUC *usecase.AccessToken, refreshTokenUC *usecase.RefreshToken) *AuthHandler {
	return &AuthHandler{
		accessTokenUC:  loginUserUC,
		refreshTokenUC: refreshTokenUC,
	}
}

func (h *AuthHandler) AccessToken(w http.ResponseWriter, r *http.Request) {
	var req dto.AccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	tokenResp, err := h.accessTokenUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, tokenResp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err)
		return
	}

	tokenResp, err := h.refreshTokenUC.Execute(r.Context(), req)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, tokenResp)
}
