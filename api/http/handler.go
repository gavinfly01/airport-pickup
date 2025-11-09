package http

import (
	"github.com/gin-gonic/gin"

	"github.com/gavin/airport-pickup/internal/app/dto"
)

type CreatePassengerReq struct {
	Name string `json:"name"`
}

type CreateDriverReq struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

func (h *Handler) createPassenger(c *gin.Context) {
	var in CreatePassengerReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := h.orderApp.CreatePassenger(in.Name)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id})
}

func (h *Handler) createDriver(c *gin.Context) {
	var in CreateDriverReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := h.orderApp.CreateDriver(in.Name, in.Rating)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id})
}

func (h *Handler) createPickupRequest(c *gin.Context) {
	var in dto.CreatePickupRequestInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := h.orderApp.CreatePickupRequest(in)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id})
}

func (h *Handler) createDriverOffer(c *gin.Context) {
	var in dto.CreateDriverOfferInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := h.orderApp.CreateDriverOffer(in)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id})
}

func (h *Handler) listBookings(c *gin.Context) {
	list, err := h.orderApp.ListBookings()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) completeBooking(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}
	if err := h.orderApp.CompleteBooking(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
