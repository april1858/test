package handler

import (
	"net/http"
	"strconv"

	"github.com/april1858/test/internal/domain"
	"github.com/april1858/test/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SubscriptionHandler связывает HTTP-слой со слоем бизнес-логики.
type SubscriptionHandler struct {
	svc domain.SubscriptionService
}

// NewSubscriptionHandler конструктор хэндлера.
func NewSubscriptionHandler(svc domain.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

// CreateSubscription
// @Summary      Создать подписку
// @Description  Добавляет новую запись о подписке пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body domain.CreateSubscriptionReq true "Данные подписки"
// @Success      201 {object} domain.Subscription
// @Failure      400 {object} map[string]string
// @Router       /api/v1/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req domain.CreateSubscriptionReq
	// Слой валидации (Пункт 4 фидбека): ShouldBindJSON автоматически триггерит валидатор go-playground с кастомными тегами
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Ошибка валидации при создании подписки", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		logger.Log.Error("Не удалось создать подписку", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetSubscription
// @Summary      Получить подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "UUID подписки"
// @Success      200 {object} domain.Subscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	sub, err := h.svc.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// UpdateSubscription
// @Summary      Обновить подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "UUID подписки"
// @Param        input body domain.UpdateSubscriptionReq true "Новые данные подписки"
// @Success      200 {object} domain.Subscription
// @Router       /api/v1/subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	var req domain.UpdateSubscriptionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.UpdateSubscription(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription
// @Summary      Удалить подписку по ID
// @Tags         subscriptions
// @Param        id path string true "UUID подписки"
// @Success      204 "No Content"
// @Router       /api/v1/subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	if err := h.svc.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSubscriptions
// @Summary      Получить список подписок (Листинг)
// @Tags         subscriptions
// @Param        limit query int false "Количество записей"
// @Param        offset query int false "Смещение"
// @Success      200 {array} domain.Subscription
// @Router       /api/v1/subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	subs, err := h.svc.ListSubscriptions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// GetTotalCost
// @Summary      Подсчет суммарной стоимости за период
// @Description  Считает сумму стоимостей подписок пользователя по названию за выбранный временной интервал
// @Tags         analytics
// @Param        user_id query string true "UUID пользователя"
// @Param        service_name query string true "Название сервиса подписки"
// @Param        from query string true "Начало периода (MM-YYYY)"
// @Param        to query string true "Конец периода (MM-YYYY)"
// @Success      200 {object} map[string]int "Пример: {"total_cost": 1200}"
// @Router       /api/v1/subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Обязательная валидация query-параметров агрегации
	if userIDStr == "" || serviceName == "" || fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters: user_id, service_name, from, to"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID format"})
		return
	}

	total, err := h.svc.CalculateTotalCost(c.Request.Context(), userID, serviceName, fromStr, toStr)
	if err != nil {
		logger.Log.Warn("Ошибка при расчете агрегации", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
