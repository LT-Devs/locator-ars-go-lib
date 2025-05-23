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
- Методы для прямой проверки доступа в условных выражениях

## Использование

### Базовый пример

```go
package main

import (
	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

func main() {
	r := gin.Default()

	// Создаем middleware с конфигурацией по умолчанию
	arsMiddleware := locatorars.NewMiddleware(locatorars.DefaultConfig())

	// Защищаем маршрут требованием права "viewallreports"
	r.GET("/reports", arsMiddleware.RequireAction("viewallreports"), reportHandler)

	r.Run(":8080")
}

func reportHandler(c *gin.Context) {
	// Обработчик выполнится только если право доступа подтверждено
	// ...
}
```

### Пользовательская конфигурация

```go
package main

import (
	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
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

### Прямая проверка доступа в условных выражениях

```go
func someHandler(c *gin.Context) {
	// Создаем middleware
	arsMiddleware := locatorars.NewMiddleware(locatorars.DefaultConfig())

	// Вариант 1: Проверка с передачей параметров вручную
	jwt := c.GetHeader("X-Authentik-Jwt")
	application := c.GetHeader("Application")
	if arsMiddleware.CheckAccess("viewreports", jwt, application) {
		// Выполняем действия, требующие права "viewreports"
		showReports(c)
	} else {
		// Обрабатываем отсутствие права доступа
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to view reports"})
	}

	// Вариант 2: Проверка с автоматическим извлечением JWT и Application из контекста
	if arsMiddleware.CheckAccessFromContext(c, "editreport") {
		// Выполняем действия, требующие права "editreport"
		editReport(c)
	} else {
		// Обрабатываем отсутствие права доступа
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to edit report"})
	}

	// Использование в условных переходах
	reportType := "standard"
	if arsMiddleware.CheckAccessFromContext(c, "viewsecretreports") {
		reportType = "secret"
	}

	// Используем reportType дальше...
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

## Методы

| Метод                                                        | Описание                                                      |
| ------------------------------------------------------------ | ------------------------------------------------------------- |
| `NewMiddleware(config Config) *Middleware`                   | Создает новый экземпляр middleware                            |
| `RequireAction(action string) gin.HandlerFunc`               | Создает middleware для защиты маршрута                        |
| `CheckAccess(action, jwt, application string) bool`          | Проверяет права доступа напрямую                              |
| `CheckAccessFromContext(c *gin.Context, action string) bool` | Проверяет права доступа, извлекая данные из контекста запроса |

## Лицензия

MIT
