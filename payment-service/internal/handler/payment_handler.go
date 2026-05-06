package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pcdoma/payment-service/internal/domain"
	"github.com/pcdoma/payment-service/internal/service"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) InternalCreate(c *gin.Context) {
	var req domain.InternalCreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// List godoc
// @Summary      List my payments
// @Description  Returns the authenticated user's payment history
// @Tags         payments
// @Produce      json
// @Param        limit   query     int  false  "Page size"   default(20)
// @Param        offset  query     int  false  "Page offset" default(0)
// @Success      200     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /payments [get]
func (h *PaymentHandler) List(c *gin.Context) {
	userID, err := uuid.Parse(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	payments, total, err := h.svc.ListByUser(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": payments, "total": total})
}

// GetByID godoc
// @Summary      Get payment by ID
// @Tags         payments
// @Produce      json
// @Param        id   path      string  true  "Payment ID (UUID)"
// @Success      200  {object}  domain.Payment
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Process godoc
// @Summary      Process a pending payment
// @Description  Mock gateway resolves the payment to completed (~95%) or failed
// @Tags         payments
// @Produce      json
// @Param        id   path      string  true  "Payment ID (UUID)"
// @Success      200  {object}  domain.Payment
// @Failure      409  {object}  map[string]string  "already processed"
// @Security     BearerAuth
// @Router       /payments/{id}/process [post]
func (h *PaymentHandler) Process(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))

	p, err := h.svc.Process(c.Request.Context(), id, userID)
	if err != nil {
		if err == service.ErrAlreadyProcessed {
			c.JSON(http.StatusConflict, gin.H{"error": "payment already processed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Refund godoc
// @Summary      Refund a payment
// @Tags         payments
// @Produce      json
// @Param        id   path      string  true  "Payment ID (UUID)"
// @Success      200  {object}  domain.Payment
// @Security     BearerAuth
// @Router       /payments/{id}/refund [post]
func (h *PaymentHandler) Refund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	p, err := h.svc.Refund(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refund"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// AdminList godoc
// @Summary      List all payments (admin)
// @Tags         admin
// @Produce      json
// @Param        limit   query     int  false  "Page size"   default(50)
// @Param        offset  query     int  false  "Page offset" default(0)
// @Success      200     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/payments [get]
func (h *PaymentHandler) AdminList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	payments, total, err := h.svc.ListAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": payments, "total": total})
}
