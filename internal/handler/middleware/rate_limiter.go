package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter возвращает Middleware для ограничения входящего трафика.
func RateLimiter(rps float64, burst int) gin.HandlerFunc { // Инициализируем лимитер: rps — запросов в секунду, burst — максимальный всплеск (емкость бакета)
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(c *gin.Context) { // Проверяем, укладывается ли запрос в лимиты
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort() // Прерываем цепочку обработки запроса
			return
		}
		c.Next() // Передаем управление следующему хэндлеру
	}

}
