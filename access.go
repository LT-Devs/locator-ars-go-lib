package locatorars

import (
	"fmt"
	"net/http"
	"time"
)

// AccessResponse представляет ответ от сервиса проверки прав доступа
type AccessResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
}

// AccessClient клиент для проверки прав доступа
type AccessClient struct {
	config  Config
	client  *http.Client
	logger  Logger
}

// NewAccessClient создает новый клиент для проверки прав доступа
func NewAccessClient(config Config) *AccessClient {
	var logger Logger
	if config.Logger != nil {
		logger = config.Logger
	} else {
		logger = NewDefaultLogger(config.LogLevel)
	}
	
	return &AccessClient{
		config: config,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

// CheckAccess проверяет права доступа для указанного действия
func (ac *AccessClient) CheckAccess(action, jwt, application string) (bool, error) {
	startTime := time.Now()
	
	// Формируем URL запроса
	url := fmt.Sprintf("%s?action=%s", ac.config.URL, action)
	ac.logger.Debug("Making access check request: URL=%s, Action=%s", url, action)
	
	// Создаем HTTP запрос
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ac.logger.Error("Failed to create request: %v", err)
		if ac.config.AllowOnFailure {
			return true, err
		}
		return false, err
	}
	
	// Добавляем необходимые заголовки
	req.Header.Set("X-Authentik-Jwt", jwt)
	req.Header.Set("Application", application)
	ac.logger.Debug("Request headers set: Application=%s, JWT present=%v", application, len(jwt) > 0)
	
	// Выполняем запрос
	ac.logger.Debug("Sending access check request...")
	resp, err := ac.client.Do(req)
	if err != nil {
		ac.logger.Error("HTTP request failed: %v", err)
		if ac.config.AllowOnFailure {
			ac.logger.Info("Access allowed on failure due to configuration")
			return true, err
		}
		return false, err
	}
	defer resp.Body.Close()
	
	elapsedMs := time.Since(startTime).Milliseconds()
	ac.logger.Debug("Access check response received in %d ms: StatusCode=%d", elapsedMs, resp.StatusCode)
	
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		ac.logger.Error("Access service returned non-200 status: %d", resp.StatusCode)
		if ac.config.AllowOnFailure {
			ac.logger.Info("Access allowed on failure due to configuration")
			return true, fmt.Errorf("access service returned non-200 status: %d", resp.StatusCode)
		}
		return false, fmt.Errorf("access service returned non-200 status: %d", resp.StatusCode)
	}
	
	// Если статус OK, считаем что доступ разрешен
	ac.logger.Debug("Access check successful, access granted")
	return true, nil
} 