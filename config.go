package locatorars

import (
	"log"
	"os"
)

// LogLevel определяет уровень логирования
type LogLevel int

const (
	// LogLevelNone - отключает логирование
	LogLevelNone LogLevel = iota
	// LogLevelError - логирует только ошибки
	LogLevelError
	// LogLevelInfo - логирует ошибки и информационные сообщения
	LogLevelInfo
	// LogLevelDebug - логирует всё, включая отладочную информацию
	LogLevelDebug
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// DefaultLogger реализация логгера по умолчанию
type DefaultLogger struct {
	logger *log.Logger
	level  LogLevel
}

// Debug логирует отладочное сообщение
func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	if l.level >= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info логирует информационное сообщение
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	if l.level >= LogLevelInfo {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

// Error логирует сообщение об ошибке
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	if l.level >= LogLevelError {
		l.logger.Printf("[ERROR] "+format, args...)
	}
}

// Config определяет конфигурацию для middleware проверки прав доступа
type Config struct {
	// URL сервиса locator-ars для проверки прав доступа
	// По умолчанию: "http://locator/api/v1/ars/check"
	URL string

	// Политика действий в случае недоступности сервиса проверки прав
	// true - разрешить доступ если сервис недоступен, 
	// false - запретить доступ если сервис недоступен
	AllowOnFailure bool
	
	// Уровень логирования
	LogLevel LogLevel
	
	// Пользовательский логгер (если nil, будет использован логгер по умолчанию)
	Logger Logger
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		URL:            "http://locator/api/v1/ars/check",
		AllowOnFailure: false,
		LogLevel:       LogLevelError,
		Logger:         nil,
	}
}

// NewDefaultLogger создает логгер по умолчанию с указанным уровнем
func NewDefaultLogger(level LogLevel) Logger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
} 