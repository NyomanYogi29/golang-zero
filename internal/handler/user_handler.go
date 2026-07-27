package handler

import (
	"encoding/json"
	"latihan/internal/global"
	"latihan/internal/model"
	"latihan/internal/service"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserService(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := global.NewErrorResponse("Method tidak diizinkan", "METHOD_NOT_ALLOWED")
		global.WriteJSON(w, http.StatusMethodNotAllowed, res)
		return
	}

	var req model.UserRegisterRequestSchema
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := global.NewErrorResponse("Payload request tidak valid", "BAD_REQUEST")
		global.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	user, err := h.userService.Register(r.Context(), req)
	if err != nil {
		res := global.NewErrorResponse(err.Error(), "REGISTER_FAILED")
		global.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	user.Password = ""

	res := global.NewSuccessResponse("Berhasil mendaftarkan akun", user)
	global.WriteJSON(w, http.StatusCreated, res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := global.NewErrorResponse("Method is not allowed", "METHOD_NOT_ALLOWED")
		global.WriteJSON(w, http.StatusMethodNotAllowed, res)
		return
	}

	var req model.UserLoginRequestSchema
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := global.NewErrorResponse("Invalid request payload", "BAD_REQUEST")
		global.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	loginResult, err := h.userService.Login(r.Context(), req)
	if err != nil {
		res := global.NewErrorResponse(err.Error(), "LOGIN_FAILED")
		global.WriteJSON(w, http.StatusUnauthorized, res)
		return
	}

	res := global.NewSuccessResponse("Login successful", loginResult)
	global.WriteJSON(w, http.StatusOK, res)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		res := global.NewErrorResponse("User ID does not found", "UNAUTHORIZED")
		global.WriteJSON(w, http.StatusUnauthorized, res)
		return
	}

	if err := h.userService.Logout(r.Context(), userID); err != nil {
		res := global.NewErrorResponse("Error processing logout", "INTERNAL_ERROR")
		global.WriteJSON(w, http.StatusInternalServerError, res)
		return
	}

	res := global.NewSuccessResponse("Logout successfully", nil)
	global.WriteJSON(w, http.StatusOK, res)
}
