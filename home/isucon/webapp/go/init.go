package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func routePostInitialize(c *gin.Context) {
	rd.FlushDb()
	os.RemoveAll("/home/isucon/webapp/public/slots/")
	c.String(200, "OK")
}
