package main

import (
	"os"
)

func routePostInitialize() (int, string) {
	keys, _ := rd.Keys("isu4:*").Result()
	for i := range keys {
		key := keys[i]
		rd.Del(key)
	}
	path := getDir("log")
	os.RemoveAll(path)

	os.RemoveAll("/home/isucon/webapp/public/slots/")

	return 200, "OK"
}
