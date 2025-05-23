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

	// Пример обработчика с условной проверкой прав доступа
	r.GET("/conditionalreport", func(c *gin.Context) {
		// Вариант 1: Проверка доступа извлекая данные из контекста
		if arsMiddleware.CheckAccessFromContext(c, "viewdetailedreport") {
			// Пользователь имеет право на просмотр подробного отчета
			c.JSON(http.StatusOK, gin.H{
				"report": "Подробный отчет с секретными данными",
				"access_level": "detailed",
			})
		} else {
			// Пользователь не имеет права, показываем упрощенный отчет
			c.JSON(http.StatusOK, gin.H{
				"report": "Упрощенный отчет без секретных данных",
				"access_level": "basic",
			})
		}
	})

	// Пример с проверкой нескольких прав доступа
	r.GET("/admin/dashboard", func(c *gin.Context) {
		// Получаем заголовки вручную
		jwt := c.GetHeader("X-Authentik-Jwt")
		application := c.GetHeader("Application")
		
		// Проверяем разные права
		canViewDashboard := arsMiddleware.CheckAccess("viewdashboard", jwt, application)
		canManageUsers := arsMiddleware.CheckAccess("manageusers", jwt, application)
		canExportData := arsMiddleware.CheckAccess("exportdata", jwt, application)
		
		// Формируем ответ на основе проверки прав
		c.JSON(http.StatusOK, gin.H{
			"permissions": map[string]bool{
				"view_dashboard": canViewDashboard,
				"manage_users":   canManageUsers,
				"export_data":    canExportData,
			},
			"features_available": []string{
				"base_dashboard",
				canManageUsers && "user_management",
				canExportData && "data_export",
			},
		})
	})

	// Запуск сервера
	r.Run(":8080")
} 