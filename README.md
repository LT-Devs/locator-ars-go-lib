# locator-ars-go-lib

Библиотека Go для интеграции с системой управления доступом Locator ARS в приложениях на базе Gin framework.

## Требования

- Go 1.24 или новее
- Gin Framework

## Установка

```bash
go get github.com/kikowoll/locator-ars-go-lib
```

## Описание

`locator-ars-go-lib` - это библиотека, предоставляющая middleware для Gin framework, которая выполняет проверку прав доступа перед выполнением обработчиков запросов.

Middleware выполняет запрос к сервису Locator ARS для проверки наличия у пользователя прав на выполнение определенного действия.

## Особенности

- Интеграция с Gin framework
- Настраиваемый URL сервиса проверки прав доступа
- Возможность разрешить или запретить доступ при недоступности сервиса проверки
- Передача JWT токена и идентификатора приложения

## Использование

### Базовый пример

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-gin-postgres/database"
	"github.com/yourusername/go-gin-postgres/handlers"
	"github.com/yourusername/go-gin-postgres/middleware"

	// Импортируем нашу библиотеку
	locatorars "github.com/kikowoll/locator-ars-go-lib"
)

func SetupRoutes(router *gin.Engine, db *database.DBConnections) {
	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	routerApi := router.Group("/api/v1")

	authMiddleware := middleware.AuthMiddleware()

	// Создаем middleware для проверки прав доступа
	arsMiddleware := locatorars.NewMiddleware(locatorars.DefaultConfig())

	petitionGroup := routerApi.Group("/petitions")
	petitionGroup.Use(authMiddleware)
	{
		// Вариант 1: Применение к конкретному маршруту
		petitionGroup.GET("", arsMiddleware.RequireAction("viewpetitions"), handlers.GetPetitions)
		petitionGroup.GET("/unregistered", arsMiddleware.RequireAction("viewunregistered"), handlers.GetUnregisteredPetitions)

		// Остальные маршруты...
		petitionGroup.GET("/:id", handlers.GetPetition)

		// Вариант 2: Применение к маршруту, требующему повышенных привилегий
		petitionGroup.POST("", arsMiddleware.RequireAction("createpetition"), handlers.CreatePetition)
		petitionGroup.POST("/create-multiple", arsMiddleware.RequireAction("createpetitions"), handlers.CreatePetitions)

		// Остальные маршруты...
	}

	// Вариант 3: Создание новой группы со своим middleware для проверки прав
	adminGroup := routerApi.Group("/admin")
	adminGroup.Use(authMiddleware)
	adminGroup.Use(arsMiddleware.RequireAction("adminaccess"))
	{
		// Все маршруты в этой группе требуют действия "adminaccess"
		adminGroup.GET("/stats", handlers.GetStats)
		adminGroup.POST("/settings", handlers.UpdateSettings)
	}
}
```

### Пользовательская конфигурация

```go
package main

import (
	"github.com/gin-gonic/gin"
	locatorars "github.com/kikowoll/locator-ars-go-lib"
)

func main() {
	r := gin.Default()

	// Создаем middleware с пользовательской конфигурацией
	config := locatorars.Config{
		URL:            "http://custom-locator/api/v2/ars/check",
		AllowOnFailure: true, // Разрешаем доступ при недоступности сервиса проверки
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Защищаем маршрут
	r.GET("/admin", arsMiddleware.RequireAction("adminaccess"), adminHandler)

	r.Run(":8080")
}

func adminHandler(c *gin.Context) {
	// ...
}
```

## Параметры конфигурации

| Параметр       | Тип    | По умолчанию                      | Описание                                                         |
| -------------- | ------ | --------------------------------- | ---------------------------------------------------------------- |
| URL            | string | "http://locator/api/v1/ars/check" | URL сервиса проверки прав доступа                                |
| AllowOnFailure | bool   | false                             | Политика доступа при ошибке: true - разрешить, false - запретить |

## Требуемые HTTP заголовки

Для корректной работы middleware клиент должен передавать следующие HTTP заголовки:

- `X-Authentik-Jwt`: JWT токен пользователя
- `Application`: Идентификатор приложения

## Коды ответов

Middleware может возвращать следующие HTTP статусы:

- `401 Unauthorized`: Отсутствует заголовок X-Authentik-Jwt
- `400 Bad Request`: Отсутствует заголовок Application
- `403 Forbidden`: Доступ запрещен
- `500 Internal Server Error`: Ошибка при проверке доступа (если AllowOnFailure=false)

## Лицензия

MIT
