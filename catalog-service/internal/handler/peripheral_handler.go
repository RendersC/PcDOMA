package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pcdoma/catalog-service/internal/domain"
	mongorepo "github.com/pcdoma/catalog-service/internal/repository/mongo"
)

type PeripheralHandler struct {
	repo *mongorepo.PeripheralRepository
}

func NewPeripheralHandler(repo *mongorepo.PeripheralRepository) *PeripheralHandler {
	return &PeripheralHandler{repo: repo}
}

func (h *PeripheralHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	peripherals, total, err := h.repo.List(c.Request.Context(), c.Query("type"), c.Query("location_id"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list peripherals"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peripherals, "total": total})
}

func (h *PeripheralHandler) ListByLocation(c *gin.Context) {
	pcID := c.Param("id")
	_ = pcID // In a real scenario we'd look up the PC's locationID
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	locationID := c.Query("location_id")

	peripherals, total, err := h.repo.List(c.Request.Context(), "", locationID, limit, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list peripherals"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peripherals, "total": total})
}

func (h *PeripheralHandler) GetByID(c *gin.Context) {
	p, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "peripheral not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PeripheralHandler) Create(c *gin.Context) {
	var req domain.CreatePeripheralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &domain.Peripheral{
		Name:         req.Name,
		Type:         req.Type,
		Description:  req.Description,
		Specs:        req.Specs,
		PricePerHour: req.PricePerHour,
		PricePerDay:  req.PricePerDay,
		LocationID:   req.LocationID,
		Images:       req.Images,
	}
	if err := h.repo.Create(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create peripheral"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *PeripheralHandler) Update(c *gin.Context) {
	p, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "peripheral not found"})
		return
	}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PeripheralHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete peripheral"})
		return
	}
	c.Status(http.StatusNoContent)
}
