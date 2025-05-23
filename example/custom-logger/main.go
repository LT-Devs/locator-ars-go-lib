package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	locatorars "github.com/LT-Devs/locator-ars-go-lib"
)

// CustomLogger пример пользовательского логгера
type CustomLogger struct {
	file *os.File
}

// Debug логирует отладочное сообщение
func (l *CustomLogger) Debug(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

// Info логирует информационное сообщение
func (l *CustomLogger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Error логирует сообщение об ошибке
func (l *CustomLogger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

// log записывает отформатированное сообщение с указанным уровнем
func (l *CustomLogger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
	
	// Вывод в консоль
	fmt.Print(logLine)
	
	// Запись в файл, если файл открыт
	if l.file != nil {
		l.file.WriteString(logLine)
	}
}

func main() {
	// Открываем файл для логов
	logFile, err := os.OpenFile("access_check.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()
	
	// Создаем кастомный логгер
	customLogger := &CustomLogger{
		file: logFile,
	}
	
	// Создаем экземпляр gin
	r := gin.Default()

	// Создаем middleware с пользовательским логгером
	config := locatorars.Config{
		URL:            "http://locator/api/v1/ars/check",
		AllowOnFailure: false,
		Logger:         customLogger, // Передаем пользовательский логгер
	}
	arsMiddleware := locatorars.NewMiddleware(config)

	// Пример маршрута
	r.GET("/reports", arsMiddleware.RequireAction("viewallreports"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"reports": []string{"report1", "report2", "report3"},
		})
	})

	// Запуск сервера
	customLogger.Info("Starting server on :8080")
	r.Run(":8080")
} 