package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger инициализирует глобальный логер в зависимости от окружения (development/production).
func InitLogger(env string) {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// Понятный формат времени
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel =
			zapcore.CapitalColorLevelEncoder // Цветной вывод в консоль разработчика
	}
	var err error
	Log, err = config.Build()
	if err != nil {
		panic("Не удалось запустить логер: " + err.Error())
	}
}
