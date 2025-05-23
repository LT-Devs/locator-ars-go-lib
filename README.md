# locator-ars-go-lib

Библиотека Go для интеграции с системой управления доступом Locator ARS в приложениях на базе Gin framework.

## Требования

- Go 1.24 или новее
- Gin Framework

## Установка

```bash
go get github.com/LT-Devs/locator-ars-go-lib
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
- Гибкая система логирования с возможностью кастомизации

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

### Настройка логирования

```go
package main

import (
	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

func main() {
	r := gin.Default()

	// Создаем middleware с включенным подробным логированием
	config := locatorars.Config{
		URL:            "http://locator/api/v1/ars/check",
		AllowOnFailure: false,
		LogLevel:       locatorars.LogLevelDebug, // Подробное логирование
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Защищаем маршрут
	r.GET("/reports", arsMiddleware.RequireAction("viewreports"), handleReports)

	// Динамическое изменение уровня логирования
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

		c.JSON(http.StatusOK, gin.H{"message": "Log level changed"})
	})

	r.Run(":8080")
}
```

### Пользовательский логгер

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

// CustomLogger пользовательский логгер, реализующий интерфейс locatorars.Logger
type CustomLogger struct {
	file *os.File
}

func (l *CustomLogger) Debug(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

func (l *CustomLogger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *CustomLogger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

func (l *CustomLogger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Вывод в консоль и файл
	fmt.Print(logLine)
	if l.file != nil {
		l.file.WriteString(logLine)
	}
}

func main() {
	// Создаем файл для логов
	logFile, _ := os.OpenFile("access_check.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()

	// Создаем экземпляр кастомного логгера
	customLogger := &CustomLogger{file: logFile}

	// Используем кастомный логгер
	config := locatorars.Config{
		URL:    "http://locator/api/v1/ars/check",
		Logger: customLogger,
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Далее настраиваем маршруты...
}
```

## Параметры конфигурации

| Параметр       | Тип      | По умолчанию                      | Описание                                                                |
| -------------- | -------- | --------------------------------- | ----------------------------------------------------------------------- |
| URL            | string   | "http://locator/api/v1/ars/check" | URL сервиса проверки прав доступа                                       |
| AllowOnFailure | bool     | false                             | Политика доступа при ошибке: true - разрешить, false - запретить        |
| LogLevel       | LogLevel | LogLevelError                     | Уровень логирования при использовании стандартного логгера              |
| Logger         | Logger   | nil                               | Пользовательский логгер, если nil, будет использован стандартный логгер |

## Уровни логирования

| Уровень       | Описание                                    |
| ------------- | ------------------------------------------- |
| LogLevelNone  | Отключает логирование                       |
| LogLevelError | Логирует только ошибки                      |
| LogLevelInfo  | Логирует ошибки и информационные сообщения  |
| LogLevelDebug | Логирует всё, включая отладочную информацию |

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
| `SetLogLevel(level LogLevel)`                                | Устанавливает уровень логирования для стандартного логгера    |

## Интерфейс Logger

Чтобы создать пользовательский логгер, реализуйте следующий интерфейс:

```go
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}
```

## Лицензия

MIT
