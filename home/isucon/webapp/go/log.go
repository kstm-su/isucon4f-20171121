package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v2"
)

func getLogPath(advrId string) string {
	dir := getDir("log")
	splitted := strings.Split(advrId, "/")
	return dir + "/" + splitted[len(splitted)-1]
}

func getLog(id string) map[string][]ClickLog {
	path := getLogPath(id)
	result := map[string][]ClickLog{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return result
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_SH)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\n")
		sp := strings.Split(line, "\t")
		ad_id := sp[0]
		user := sp[1]
		agent := sp[2]
		if agent == "" {
			agent = "unknown"
		}
		gender, age := decodeUserKey(sp[1])
		if result[ad_id] == nil {
			result[ad_id] = []ClickLog{}
		}
		data := ClickLog{ad_id, user, agent, gender, age}
		result[ad_id] = append(result[ad_id], data)
	}

	return result
}

