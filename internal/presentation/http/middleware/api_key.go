package middleware

import (
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"net/http"
	"strings"

	"bufio"
	"database/sql"
	"os"
	// Postgres driver - adjust if needed
)

type APIKeyMiddleware struct {
	config *config.Config
}

func NewAPIKeyMiddleware(cfg *config.Config) *APIKeyMiddleware {
	return &APIKeyMiddleware{config: cfg}
}

func (m *APIKeyMiddleware) ValidateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key") // Или откуда вы хотите его получать
		if apiKey == "" {
			logging.Logger.Warn("API Key отсутствует в запросе")
			http.Error(w, "API Key required", http.StatusUnauthorized)
			return
		}

		// Проверяем API Key в зависимости от источника (static, file, db)
		valid := false
		switch strings.ToLower(m.config.APIKeys.Source) {
		case "static":
			for _, key := range m.config.APIKeys.Keys {
				if apiKey == key {
					valid = true
					break
				}
			}
		case "file":
			valid = m.validateAPIKeyFromFile(apiKey)
		case "db":
			valid = m.validateAPIKeyFromDB(apiKey)
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

		// API Key валиден, продолжаем обработку
		next.ServeHTTP(w, r)
	})
}

func (m *APIKeyMiddleware) validateAPIKeyFromFile(apiKey string) bool {
	filePath := m.config.APIKeys.FilePath
	file, err := os.Open(filePath)
	if err != nil {
		logging.Logger.Error("Ошибка открытия файла с API ключами", "error", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if apiKey == key {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		logging.Logger.Error("Ошибка чтения файла с API ключами", "error", err)
		return false
	}
	return false
}

func (m *APIKeyMiddleware) validateAPIKeyFromDB(apiKey string) bool {
	dbConnString := m.config.APIKeys.DBConnString
	db, err := sql.Open("postgres", dbConnString) // Update if needed
	if err != nil {
		logging.Logger.Error("Ошибка подключения к базе данных", "error", err)
		return false
	}
	defer db.Close()

	var foundKey string
	err = db.QueryRow("SELECT api_key FROM api_keys WHERE api_key = $1 AND enabled = true", apiKey).Scan(&foundKey)

	if err == sql.ErrNoRows {
		logging.Logger.Warn("API Key не найден в базе данных")
		return false
	} else if err != nil {
		logging.Logger.Error("Ошибка запроса к базе данных", "error", err)
		return false
	}

	return foundKey == apiKey
}
