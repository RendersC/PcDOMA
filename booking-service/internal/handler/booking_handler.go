package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pcdoma/booking-service/internal/domain"
	"github.com/pcdoma/booking-service/internal/service"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// Create godoc
// @Summary      Create a booking
// @Description  Creates a custom (PC + peripherals) or setup booking, takes payment, and routes it to a worker
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreateBookingRequest  true  "Booking request"
// @Success      201      {object}  domain.Booking
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string  "PC not available"
// @Security     BearerAuth
// @Router       /bookings [post]
func (h *BookingHandler) Create(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	booking, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		if err == service.ErrPCNotAvailable {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

// GetByID godoc
// @Summary      Get booking by ID
// @Tags         bookings
// @Produce      json
// @Param        id   path      string  true  "Booking ID (UUID)"
// @Success      200  {object}  domain.Booking
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /bookings/{id} [get]
func (h *BookingHandler) GetByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	role := c.GetHeader("X-User-Role")

	b, err := h.svc.GetByID(c.Request.Context(), id, userID, role)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

// ListByUser godoc
// @Summary      List my bookings
// @Tags         bookings
// @Produce      json
// @Param        status  query     string  false  "Filter by status"
// @Param        limit   query     int     false  "Page size"   default(20)
// @Param        offset  query     int     false  "Page offset" default(0)
// @Success      200     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /bookings [get]
func (h *BookingHandler) ListByUser(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, total, err := h.svc.ListByUser(c.Request.Context(), userID, c.Query("status"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings, "total": total})
}

// Cancel godoc
// @Summary      Cancel a booking
// @Tags         bookings
// @Produce      json
// @Param        id   path      string  true  "Booking ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      409  {object}  map[string]string  "cannot cancel"
// @Security     BearerAuth
// @Router       /bookings/{id} [delete]
func (h *BookingHandler) Cancel(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	role := c.GetHeader("X-User-Role")

	if err := h.svc.Cancel(c.Request.Context(), id, userID, role); err != nil {
		switch err {
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case service.ErrInvalidTransition:
			c.JSON(http.StatusConflict, gin.H{"error": "cannot cancel this booking"})
		default:
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

func (h *BookingHandler) AdminList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, total, err := h.svc.ListAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings, "total": total})
}

// Worker handlers

func (h *BookingHandler) WorkerList(c *gin.Context) {
	locationID := c.Query("location_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, total, err := h.svc.ListByLocation(c.Request.Context(), locationID, "pending_worker", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings, "total": total})
}

func (h *BookingHandler) WorkerActiveList(c *gin.Context) {
	locationID := c.Query("location_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	bookings, total, err := h.svc.ListByLocation(c.Request.Context(), locationID, "active", limit, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings, "total": total})
}

func (h *BookingHandler) WorkerHistory(c *gin.Context) {
	locationID := c.Query("location_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	bookings, total, err := h.svc.ListByLocation(c.Request.Context(), locationID, "completed", limit, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings, "total": total})
}

// WorkerAccept godoc
// @Summary      Worker accepts a pending booking
// @Tags         worker
// @Produce      json
// @Param        id   path      string  true  "Booking ID (UUID)"
// @Success      200  {object}  domain.Booking
// @Failure      409  {object}  map[string]string  "invalid transition"
// @Security     BearerAuth
// @Router       /worker/bookings/{id}/accept [patch]
func (h *BookingHandler) WorkerAccept(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	workerID, _ := uuid.Parse(c.GetHeader("X-User-ID"))

	b, err := h.svc.WorkerAccept(c.Request.Context(), id, workerID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) WorkerReject(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	workerID, _ := uuid.Parse(c.GetHeader("X-User-ID"))

	var req domain.WorkerActionRequest
	c.ShouldBindJSON(&req)

	b, err := h.svc.WorkerReject(c.Request.Context(), id, workerID, req.RejectionReason)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) WorkerDeliver(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	b, err := h.svc.WorkerDeliver(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) WorkerComplete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	b, err := h.svc.WorkerComplete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}
