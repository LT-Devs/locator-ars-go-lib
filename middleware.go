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
		entitlements := c.GetHeader("X-Authentik-Entitlements")
		if entitlements == "" {
			m.logger.Info("Missing X-Authentik-Entitlements header in request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing X-Authentik-Entitlements header",
			})
			return
		}

		// Application header is no longer required as authentik provides entitlements for specific applications
		application := c.GetHeader("Application")
		if application == "" {
			// Use default application identifier or extract from context if needed
			application = "default"
		}

		m.logger.Debug("Headers found: Application=%s, Entitlements present=%v", application, len(entitlements) > 0)

		// Проверяем доступ
		allowed, err := m.client.CheckAccess(action, entitlements)
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

// CheckAccess проверяет права доступа для указанного действия, Entitlements и приложения
// Возвращает true если доступ разрешен, false если запрещен
// Может использоваться напрямую в условных выражениях
func (m *Middleware) CheckAccess(action, entitlements string) bool {
	m.logger.Debug("Direct check for action: %s", action)

	allowed, err := m.client.CheckAccess(action, entitlements)
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

// CheckAccessFromContext проверяет права доступа, извлекая Entitlements и приложение из gin.Context
// Удобно для использования в обработчиках
func (m *Middleware) CheckAccessFromContext(c *gin.Context, action string) bool {
	entitlements := c.GetHeader("X-Authentik-Entitlements")

	m.logger.Debug("Context check for action: %s", action)

	if entitlements == "" {
		m.logger.Info("Context access check: missing headers for action: %s", action)
		return false
	}

	return m.CheckAccess(action, entitlements)
}

// SetLogLevel устанавливает уровень логирования для middleware
func (m *Middleware) SetLogLevel(level LogLevel) {
	if defaultLogger, ok := m.logger.(*DefaultLogger); ok {
		defaultLogger.level = level
	}
}
