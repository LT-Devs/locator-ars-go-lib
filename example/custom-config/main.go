package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

func main() {
	// Создаем экземпляр gin
	r := gin.Default()

	// Создаем middleware с пользовательской конфигурацией
	config := locatorars.Config{
		URL:            "http://custom-locator/api/v2/ars/check",
		AllowOnFailure: true, // Разрешаем доступ при недоступности сервиса проверки прав
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Пример маршрута, защищенного проверкой прав доступа
	r.GET("/admin", arsMiddleware.RequireAction("adminaccess"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "admin access granted",
		})
	})

	// Запуск сервера
	r.Run(":8080")
} 