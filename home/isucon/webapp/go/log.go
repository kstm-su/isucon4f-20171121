package main

import (
	"fmt"
	"os"
	"strings"
)

func getLogPath(advrId string) string {
	dir := getDir("log")
	splitted := strings.Split(advrId, "/")
	return dir + "/" + splitted[len(splitted)-1]
}

func getLog(id string) map[string][]ClickLog {
	path := getLogPath(id)
	s, err := rd.LRange(path, 0, -1).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "getLog: ", err)
		return nil
	}

	result := map[string][]ClickLog{}
	for _, l := range s {
		sp := strings.Split(l, "\t")
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
