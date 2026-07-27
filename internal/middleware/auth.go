package middleware

import (
	"context"
	"latihan/internal/global"
	"latihan/internal/service"
	"latihan/internal/utils"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	userService *service.UserService
}

func NewAuthMiddleware(userService *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			res := global.NewErrorResponse("Header Authorization Needed", "UNAUTHORIZED")
			global.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			res := global.NewErrorResponse("Authorization Format must be 'Bearer <token>'", "UNAUTHORIZED")
			global.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			res := global.NewErrorResponse("Invalid or expiry token", "UNAUTHORIZED")
			global.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		if !m.userService.IsSessionValid(r.Context(), claims.UserID, tokenString) {
			res := global.NewErrorResponse("Session ended or have been logged out", "UNAUTHORIZED")
			global.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
