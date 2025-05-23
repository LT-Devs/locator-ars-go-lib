package locatorars

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware предоставляет функциональность проверки прав доступа
type Middleware struct {
	client *AccessClient
	config Config
}

// NewMiddleware создает новый экземпляр middleware для проверки прав доступа
func NewMiddleware(config Config) *Middleware {
	return &Middleware{
		client: NewAccessClient(config),
		config: config,
	}
}

// RequireAction создает middleware, который требует указанное действие
func (m *Middleware) RequireAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем необходимые заголовки
		jwt := c.GetHeader("X-Authentik-Jwt")
		if jwt == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing X-Authentik-Jwt header",
			})
			return
		}

		application := c.GetHeader("Application")
		if application == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing Application header",
			})
			return
		}

		// Проверяем доступ
		allowed, err := m.client.CheckAccess(action, jwt, application)
		if err != nil {
			log.Printf("Error checking access: %v", err)
			if !m.config.AllowOnFailure {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to check access",
				})
				return
			}
			// Если настроено разрешать при ошибке, продолжаем выполнение
		}

		// Если доступ запрещен, возвращаем ошибку
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}

		// Если доступ разрешен, продолжаем выполнение следующего обработчика
		c.Next()
	}
} 