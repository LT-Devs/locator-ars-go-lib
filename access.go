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
}

// NewAccessClient создает новый клиент для проверки прав доступа
func NewAccessClient(config Config) *AccessClient {
	return &AccessClient{
		config: config,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CheckAccess проверяет права доступа для указанного действия
func (ac *AccessClient) CheckAccess(action, jwt, application string) (bool, error) {
	// Формируем URL запроса
	url := fmt.Sprintf("%s?action=%s", ac.config.URL, action)
	
	// Создаем HTTP запрос
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		if ac.config.AllowOnFailure {
			return true, err
		}
		return false, err
	}
	
	// Добавляем необходимые заголовки
	req.Header.Set("X-Authentik-Jwt", jwt)
	req.Header.Set("Application", application)
	
	// Выполняем запрос
	resp, err := ac.client.Do(req)
	if err != nil {
		if ac.config.AllowOnFailure {
			return true, err
		}
		return false, err
	}
	defer resp.Body.Close()
	
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		if ac.config.AllowOnFailure {
			return true, fmt.Errorf("access service returned non-200 status: %d", resp.StatusCode)
		}
		return false, fmt.Errorf("access service returned non-200 status: %d", resp.StatusCode)
	}
	
	// Если статус OK, считаем что доступ разрешен
	return true, nil
} 