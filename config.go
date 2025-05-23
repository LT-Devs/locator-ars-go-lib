package locatorars

// Config определяет конфигурацию для middleware проверки прав доступа
type Config struct {
	// URL сервиса locator-ars для проверки прав доступа
	// По умолчанию: "http://locator/api/v1/ars/check"
	URL string

	// Политика действий в случае недоступности сервиса проверки прав
	// true - разрешить доступ если сервис недоступен, 
	// false - запретить доступ если сервис недоступен
	AllowOnFailure bool
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		URL:           "http://locator/api/v1/ars/check",
		AllowOnFailure: false,
	}
} 