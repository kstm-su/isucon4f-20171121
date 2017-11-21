package main

import (
	"github.com/gin-gonic/gin"
)

func routeGetAdCount(c *gin.Context) {
	slot := c.Param("slot")
	id := c.Param("id")
	key := adKey(slot, id)

	exists, _ := rd.Exists(key).Result()
	if !exists {
		c.JSON(404, map[string]string{"error": "not_found"})
		return
	}

	rd.HIncrBy(key, "impressions", 1).Result()
	c.Status(204)
}
