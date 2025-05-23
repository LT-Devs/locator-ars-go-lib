package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

func main() {
	// Создаем экземпляр gin
	r := gin.Default()

	// Создаем middleware с конфигурацией по умолчанию
	arsMiddleware := locatorars.NewMiddleware(locatorars.DefaultConfig())

	// Пример маршрута, защищенного проверкой прав доступа
	r.GET("/reports", arsMiddleware.RequireAction("viewallreports"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"reports": []string{"report1", "report2", "report3"},
		})
	})

	// Пример с другим действием
	r.POST("/reports", arsMiddleware.RequireAction("createreport"), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Report created",
		})
	})

	// Запуск сервера
	r.Run(":8080")
} 