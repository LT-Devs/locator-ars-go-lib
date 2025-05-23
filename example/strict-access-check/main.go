package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

func main() {
	// Создаем экземпляр gin
	r := gin.Default()

	// Создаем middleware с конфигурацией и логированием
	config := locatorars.Config{
		URL:            "http://locator/api/v1/ars/check",
		AllowOnFailure: false,
		LogLevel:       locatorars.LogLevelDebug, // Включаем подробное логирование
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Общий маршрут, проверяющий доступ
	r.GET("/check-access", func(c *gin.Context) {
		// Получаем значение application из query-параметра
		application := c.DefaultQuery("app", "")
		if application == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Application parameter is required",
			})
			return
		}

		// Получаем JWT из заголовка
		jwt := c.GetHeader("X-Authentik-Jwt")
		if jwt == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-Authentik-Jwt header is required",
			})
			return
		}

		// Проверяем доступ для указанного приложения
		allowed := arsMiddleware.CheckAccess("view", jwt, application)

		// Выводим результат проверки
		c.JSON(http.StatusOK, gin.H{
			"application": application,
			"action":      "view",
			"allowed":     allowed,
		})
	})

	// Запуск сервера
	r.Run(":8080")
} 