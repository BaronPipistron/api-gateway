package middleware

import (
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"net/http"
)

type ApiKeyMiddleware struct {
	config *config.Config
}

func NewAPIKeyMiddleware(cfg *config.Config) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{config: cfg}
}
func (m *ApiKeyMiddleware) ValidateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			logging.Logger.Warn("API Key отсутствует в запросе")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		valid := false
		switch m.config.APIKeys.Source {
		case "static":
			for _, key := range m.config.APIKeys.Keys {
				if apiKey == key {
					valid = true
					break
				}
			}
		case "file":
			logging.Logger.Error("Чтение API ключей из файла не реализовано")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		case "db":
			logging.Logger.Error("Чтение API ключей из базы данных не реализовано")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		default:
			logging.Logger.Error("Неизвестный источник API ключей")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !valid {
			logging.Logger.Warn("Неверный API Key")
			http.Error(w, "Invalid API Key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
