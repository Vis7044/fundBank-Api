package controllers

import (
	"github.com/funcBank_Api/services"
	"github.com/gin-gonic/gin"
)

type SeedController struct {
	seedService *services.SeedService
}

func NewSeedController(seedService *services.SeedService) *SeedController {
	return &SeedController{
		seedService: seedService,
	}
}

func (sc *SeedController) SeedNAVData(c *gin.Context) {	
	path := "funds.txt"
	err := sc.seedService.SeedNAVData(c.Request.Context(), path)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to seed NAV data"})
		return
	}
	c.JSON(200, gin.H{"message": "NAV data seeded successfully"})
}

