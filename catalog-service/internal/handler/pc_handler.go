package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pcdoma/catalog-service/internal/domain"
	"github.com/pcdoma/catalog-service/internal/service"
)

type PCHandler struct {
	svc *service.PCService
}

func NewPCHandler(svc *service.PCService) *PCHandler {
	return &PCHandler{svc: svc}
}

func (h *PCHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	filter := domain.PCFilter{
		Status:   c.Query("status"),
		Location: c.Query("location"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Sort:     c.Query("sort"),
		Limit:    limit,
		Offset:   offset,
	}

	pcs, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pcs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   pcs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *PCHandler) GetByID(c *gin.Context) {
	pc, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pc not found"})
		return
	}
	c.JSON(http.StatusOK, pc)
}

func (h *PCHandler) Create(c *gin.Context) {
	var req domain.CreatePCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pc, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pc"})
		return
	}
	c.JSON(http.StatusCreated, pc)
}

func (h *PCHandler) Update(c *gin.Context) {
	var req domain.UpdatePCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pc, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pc"})
		return
	}
	c.JSON(http.StatusOK, pc)
}

func (h *PCHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status domain.PCStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), c.Param("id"), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *PCHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete pc"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PCHandler) GetLocations(c *gin.Context) {
	locations, err := h.svc.GetLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get locations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": locations})
}

func (h *PCHandler) GetAvailability(c *gin.Context) {
	resp, err := h.svc.GetAvailability(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pc not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PCHandler) CreateReview(c *gin.Context) {
	var req domain.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	if userName == "" {
		userName = "User"
	}

	review, err := h.svc.CreateReview(c.Request.Context(), c.Param("id"), userID, userName, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}
	c.JSON(http.StatusCreated, review)
}

func (h *PCHandler) ListReviews(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, total, err := h.svc.ListReviews(c.Request.Context(), c.Param("id"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reviews"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reviews, "total": total})
}

func (h *PCHandler) DeleteReview(c *gin.Context) {
	if err := h.svc.DeleteReview(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}
	c.Status(http.StatusNoContent)
}
