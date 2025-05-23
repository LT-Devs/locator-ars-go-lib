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

	// Пример использования middleware с логированием
	r.GET("/reports", arsMiddleware.RequireAction("viewallreports"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"reports": []string{"report1", "report2", "report3"},
		})
	})

	// Пример с проверкой доступа внутри обработчика
	r.GET("/conditional", func(c *gin.Context) {
		// При вызове CheckAccessFromContext будут записаны логи с информацией о проверке
		if arsMiddleware.CheckAccessFromContext(c, "editreports") {
			c.JSON(http.StatusOK, gin.H{
				"can_edit": true,
				"message": "You have edit permissions",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"can_edit": false,
				"message": "You have view-only permissions",
			})
		}
	})

	// Пример динамического изменения уровня логирования
	r.POST("/debug/loglevel/:level", func(c *gin.Context) {
		level := c.Param("level")
		
		switch level {
		case "none":
			arsMiddleware.SetLogLevel(locatorars.LogLevelNone)
		case "error":
			arsMiddleware.SetLogLevel(locatorars.LogLevelError)
		case "info":
			arsMiddleware.SetLogLevel(locatorars.LogLevelInfo)
		case "debug":
			arsMiddleware.SetLogLevel(locatorars.LogLevelDebug)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log level"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Log level changed to " + level,
		})
	})

	// Запуск сервера
	r.Run(":8080")
} 