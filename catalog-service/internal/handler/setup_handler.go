package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pcdoma/catalog-service/internal/domain"
	mongorepo "github.com/pcdoma/catalog-service/internal/repository/mongo"
)

type SetupHandler struct {
	repo           *mongorepo.SetupRepository
	pcRepo         *mongorepo.PCRepository
	peripheralRepo *mongorepo.PeripheralRepository
}

func NewSetupHandler(repo *mongorepo.SetupRepository, pcRepo *mongorepo.PCRepository, peripheralRepo *mongorepo.PeripheralRepository) *SetupHandler {
	return &SetupHandler{repo: repo, pcRepo: pcRepo, peripheralRepo: peripheralRepo}
}

func (h *SetupHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	setups, total, err := h.repo.List(c.Request.Context(), c.Query("category"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list setups"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": setups, "total": total})
}

func (h *SetupHandler) GetByID(c *gin.Context) {
	s, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || s == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setup not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SetupHandler) Create(c *gin.Context) {
	var req domain.CreateSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var totalPerHour, totalPerDay float64

	// Auto-calculate price from PC + peripherals
	if pc, err := h.pcRepo.GetByID(ctx, req.PCID); err == nil && pc != nil {
		totalPerHour += pc.PricePerHour
		totalPerDay += pc.PricePerDay
	}
	for _, pID := range req.PeripheralIDs {
		if p, err := h.peripheralRepo.GetByID(ctx, pID); err == nil && p != nil {
			totalPerHour += p.PricePerHour
			totalPerDay += p.PricePerDay
		}
	}
	if req.DiscountPercent > 0 {
		totalPerHour = totalPerHour * (1 - req.DiscountPercent/100)
		totalPerDay = totalPerDay * (1 - req.DiscountPercent/100)
	}

	setup := &domain.Setup{
		Name:              req.Name,
		Category:          req.Category,
		PCID:              req.PCID,
		PeripheralIDs:     req.PeripheralIDs,
		TotalPricePerHour: totalPerHour,
		TotalPricePerDay:  totalPerDay,
		DiscountPercent:   req.DiscountPercent,
		Description:       req.Description,
		Images:            req.Images,
	}
	if err := h.repo.Create(ctx, setup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create setup"})
		return
	}
	c.JSON(http.StatusCreated, setup)
}

func (h *SetupHandler) Update(c *gin.Context) {
	s, err := h.repo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || s == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setup not found"})
		return
	}
	if err := c.ShouldBindJSON(s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SetupHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete setup"})
		return
	}
	c.Status(http.StatusNoContent)
}
