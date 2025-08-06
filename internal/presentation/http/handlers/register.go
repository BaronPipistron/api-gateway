package handlers

import (
	"encoding/json"
	"github.com/BaronPipistron/api-gateway/internal/app/service"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"net/http"
)

type RegisterHandler struct {
	authService *service.AuthService
}

func NewRegisterHandler(authService *service.AuthService) *RegisterHandler {
	return &RegisterHandler{authService: authService}
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		logging.Logger.Error("Ошибка декодирования запроса", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.authService.Register(registerRequest)
	if err != nil {
		logging.Logger.Error("Ошибка при регистрации пользователя", "error", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
