package locatorars

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware предоставляет функциональность проверки прав доступа
type Middleware struct {
	client *AccessClient
	config Config
	logger Logger
}

// NewMiddleware создает новый экземпляр middleware для проверки прав доступа
func NewMiddleware(config Config) *Middleware {
	var logger Logger
	if config.Logger != nil {
		logger = config.Logger
	} else {
		logger = NewDefaultLogger(config.LogLevel)
	}
	
	return &Middleware{
		client: NewAccessClient(config),
		config: config,
		logger: logger,
	}
}

// RequireAction создает middleware, который требует указанное действие
func (m *Middleware) RequireAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Debug("Checking access for action: %s", action)
		
		// Получаем необходимые заголовки
		jwt := c.GetHeader("X-Authentik-Jwt")
		if jwt == "" {
			m.logger.Info("Missing X-Authentik-Jwt header in request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing X-Authentik-Jwt header",
			})
			return
		}

		application := c.GetHeader("Application")
		if application == "" {
			m.logger.Info("Missing Application header in request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing Application header",
			})
			return
		}

		m.logger.Debug("Headers found: Application=%s, JWT present=%v", application, len(jwt) > 0)

		// Проверяем доступ
		allowed, err := m.client.CheckAccess(action, jwt, application)
		if err != nil {
			m.logger.Error("Error checking access: %v", err)
			if !m.config.AllowOnFailure {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to check access",
				})
				return
			}
			m.logger.Info("Access allowed on failure due to configuration")
			// Если настроено разрешать при ошибке, продолжаем выполнение
		}

		// Если доступ запрещен, возвращаем ошибку
		if !allowed {
			m.logger.Info("Access denied for action: %s, application: %s", action, application)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}

		m.logger.Info("Access granted for action: %s, application: %s", action, application)
		// Если доступ разрешен, продолжаем выполнение следующего обработчика
		c.Next()
	}
}

// CheckAccess проверяет права доступа для указанного действия, JWT токена и приложения
// Возвращает true если доступ разрешен, false если запрещен
// Может использоваться напрямую в условных выражениях
func (m *Middleware) CheckAccess(action, jwt, application string) bool {
	m.logger.Debug("Direct check for action: %s, application: %s", action, application)
	
	allowed, err := m.client.CheckAccess(action, jwt, application)
	if err != nil {
		m.logger.Error("Error in direct access check: %v", err)
		// Возвращаем значение в соответствии с политикой обработки ошибок
		return m.config.AllowOnFailure
	}
	
	if allowed {
		m.logger.Info("Direct access check: granted for action: %s", action)
	} else {
		m.logger.Info("Direct access check: denied for action: %s", action)
	}
	
	return allowed
}

// CheckAccessFromContext проверяет права доступа, извлекая JWT токен и приложение из gin.Context
// Удобно для использования в обработчиках
func (m *Middleware) CheckAccessFromContext(c *gin.Context, action string) bool {
	jwt := c.GetHeader("X-Authentik-Jwt")
	application := c.GetHeader("Application")
	
	m.logger.Debug("Context check for action: %s, application: %s", action, application)
	
	if jwt == "" || application == "" {
		m.logger.Info("Context access check: missing headers for action: %s", action)
		return false
	}
	
	return m.CheckAccess(action, jwt, application)
}

// SetLogLevel устанавливает уровень логирования для middleware
func (m *Middleware) SetLogLevel(level LogLevel) {
	if defaultLogger, ok := m.logger.(*DefaultLogger); ok {
		defaultLogger.level = level
	}
} 