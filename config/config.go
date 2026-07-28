package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config структура хранит все конфигурационные данные приложения.
type Config struct {
	AppPort        string  `mapstructure:"APP_PORT"`
	AppEnv         string  `mapstructure:"APP_ENV"`
	DBHost         string  `mapstructure:"DB_HOST"`
	DBPort         string  `mapstructure:"DB_PORT"`
	DBUser         string  `mapstructure:"DB_USER"`
	DBPassword     string  `mapstructure:"DB_PASSWORD"`
	DBName         string  `mapstructure:"DB_NAME"`
	DBSslMode      string  `mapstructure:"DB_SSLMODE"`
	RateLimitRPS   float64 `mapstructure:"RATE_LIMIT_RPS"`   // Количество запросов в секунду для лимитера
	RateLimitBurst int     `mapstructure:"RATE_LIMIT_BURST"` // Вместимость бакета (всплеск запросов)
}

// LoadConfig читает файл .env или переменные окружения ОС.

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// ХИТРОСТЬ ДЛЯ VIPER: Явно привязываем переменные окружения,
	// чтобы Unmarshal увидел их без наличия .env файла
	viper.BindEnv("APP_PORT")
	viper.BindEnv("APP_ENV")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("RATE_LIMIT_RPS")
	viper.BindEnv("RATE_LIMIT_BURST")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Предупреждение: .env файл не найден, используются системные переменные: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
