package main

import (
	"github.com/gin-gonic/gin"
)

func routeGetAdWithId(c *gin.Context) {
	slot := c.Param("slot")
	id := c.Param("id")
	ad := getAd(c.Request, slot, id)
	if ad != nil {
		c.JSON(200, ad)
	} else {
		c.JSON(404, map[string]string{"error": "not_found"})
	}
}
