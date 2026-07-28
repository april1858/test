package middleware

import (
	"time"

	"github.com/april1858/test/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger логирует информацию о входящих HTTP-запросах.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// Обработка запроса
		c.Next()
		// Сбор метрик после выполнения запроса
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		// Записываем структурированный лог
		logger.Log.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
