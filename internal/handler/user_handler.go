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
