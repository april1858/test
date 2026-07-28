package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/april1858/test/docs"

	"github.com/april1858/test/config"
	"github.com/april1858/test/internal/handler"
	"github.com/april1858/test/internal/handler/middleware"
	"github.com/april1858/test/internal/pkg/logger"
	"github.com/april1858/test/internal/repository"
	"github.com/april1858/test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Subscription Aggregator API
// @version         1.0
// @description     REST-сервис для агрегации данных об онлайн подписках пользователей.
// @host            localhost:8080
// @BasePath        /
func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Printf("Критическая ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// 2. Инициализируем структурированный логер Zap
	logger.InitLogger(cfg.AppEnv)
	defer logger.Log.Sync() // Очистка буфера логов перед выходом
	logger.Log.Info("Конфигурация и логер успешно инициализированы")

	// 3. Подключаемся к базе данных PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSslMode)

	fmt.Println("dsn = ", dsn)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.Log.Fatal("Не удалось подключиться к СУБД PostgreSQL", zap.Error(err))
	}
	defer db.Close()
	logger.Log.Info("Успешное подключение к PostgreSQL")

	// 4. Регистрируем кастомный валидатор даты в Gin (Пункт 4 фидбека)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customdate", handler.CustomDateValidator)
		logger.Log.Info("Кастомный валидатор даты 'customdate' успешно зарегистрирован")
	}

	// 5. Инициализация слоев (Dependency Injection)
	subRepo := repository.NewPostgresRepository(db)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	// 6. Настройка HTTP-сервера Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Подключаем глобальные middleware (Пункт 4 ТЗ, Пункт 3 фидбека)
	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery()) // Защита от panic внутри хэндлеров
	r.Use(middleware.RateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst))

	// Эндпоинт для Swagger документации
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Регистрация API маршрутов (Ручки CRUDL и Агрегации)
	api := r.Group("/api/v1")
	{
		subs := api.Group("/subscriptions")
		{
			subs.POST("", subHandler.CreateSubscription)       // Create
			subs.GET("", subHandler.ListSubscriptions)         // List
			subs.GET("/:id", subHandler.GetSubscription)       // Read
			subs.PUT("/:id", subHandler.UpdateSubscription)    // Update
			subs.DELETE("/:id", subHandler.DeleteSubscription) // Delete
			subs.GET("/total", subHandler.GetTotalCost)        // Aggregation (Пункт 2 ТЗ)
		}
	}

	// 7. Реализация Graceful Shutdown (Безопасная остановка сервера без обрыва текущих запросов)
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		logger.Log.Info("Запуск HTTP сервера", zap.String("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Ошибка при работе HTTP сервера", zap.Error(err))
		}
	}()

	// Ожидаем системный сигнал завершения (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Получен сигнал завершения, завершаем работу сервера...")

	// Даем серверу 5 секунд на завершение обработки текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Сервер завершил работу некорректно (таймаут)", zap.Error(err))
	}
	logger.Log.Info("Сервер успешно и безопасно остановлен")
}
