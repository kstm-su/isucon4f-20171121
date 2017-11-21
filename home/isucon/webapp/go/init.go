package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func routePostInitialize(c *gin.Context) {
	keys, _ := rd.Keys("isu4:*").Result()
	for i := range keys {
		key := keys[i]
		rd.Del(key)
	}
	path := getDir("log")
	os.RemoveAll(path)

	os.RemoveAll("/home/isucon/webapp/public/slots/")

	c.String(200, "OK")
}
