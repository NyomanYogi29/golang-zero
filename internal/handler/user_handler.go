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

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.UserRegisterRequestSchema true "Register payload"
// @Success      201 {object} global.SuccessResponse{data=model.User}
// @Failure      400 {object} global.ErrorResponse
// @Failure      405 {object} global.ErrorResponse
// @Router       /register [post]
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

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password and receive a JWT access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.UserLoginRequestSchema true "Login payload"
// @Success      200 {object} global.SuccessResponse{data=model.UserLoginResponseSchema}
// @Failure      400 {object} global.ErrorResponse
// @Failure      401 {object} global.ErrorResponse
// @Failure      405 {object} global.ErrorResponse
// @Router       /login [post]
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

// Logout godoc
// @Summary      Logout
// @Description  Invalidate the current session for the authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} global.GlobalResponse
// @Failure      401 {object} global.ErrorResponse
// @Failure      500 {object} global.ErrorResponse
// @Router       /logout [post]
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
