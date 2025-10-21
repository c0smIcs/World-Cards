package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Config - структура для хранения конфигурационных параметров приложения
// Теги `env` используются для указания, из какой переменной окружения
// следует загружать значение поля
type Config struct {
// TgToken содержит токен для доступа к Telegram Bot API.
	// `env:"TOKEN,required"` означает, что значение будет взято из переменной TOKEN
	TgToken string `env:"TOKEN,required"`
	TranslateAPIKey string `env:"TRANSLATEAPIKEY"`
}

// Lead загружает конфигурацию из .env файла и переменных окружения
// функция сначала пытается загрузить .env файл, а затем парсит переменные
// окружения в структуре Config
func Lead() (*Config, error) {
	// пытаемся загрузить переменные из файла .env в текущей директории
	// это удобно для локальной разработки
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при загрузке файла .env: %v", err)
	}

	// создаем пустой экземпляр структуры Config
	cfg := Config{}
	// заполняем структуру cfg значениями из переменных окружения.
	// библиотека env автоматически проверит наличие обязательных полей
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга переменных: %v", err)
	}

	// возвращаем указатель на заполненную структуру конфигурации
	return &cfg, nil
}
