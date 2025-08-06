package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BaronPipistron/api-gateway/internal/app"
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"net/http"
)

type AuthService struct {
	config   *config.Config
	registry *app.UserRegistry // Add UserRegistry
}

func NewAuthService(config *config.Config, registry *app.UserRegistry) *AuthService {
	return &AuthService{config: config, registry: registry}
}

func (s *AuthService) Register(request map[string]interface{}) error {
	idpURL := s.config.IDP.RegistrationURL

	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(idpURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		logging.Logger.Error("Error registering user", "status_code", resp.StatusCode, "body", resp.Body)
		return fmt.Errorf("IDP returned error code: %d", resp.StatusCode)
	}

	username, ok := request["username"].(string)
	if !ok {
		logging.Logger.Warn("Имя пользователя отсутствует в запросе на регистрацию, не сохранено в реестре пользователей")
	} else {
		s.registry.RegisterUser(username, request)
	}
	return nil
}

func (s *AuthService) GetUser(username string) (interface{}, bool) {
	return s.registry.GetUser(username)
}
