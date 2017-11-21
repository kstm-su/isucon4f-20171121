package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func routeGetReport(c *gin.Context) {
	advrId := advertiserId(c.Request)

	if advrId == "" {
		c.Status(401)
		return
	}

	report := map[string]*Report{}
	adKeys, _ := rd.SMembers(advertiserKey(advrId)).Result()
	for _, adKey := range adKeys {
		ad, _ := rd.HGetAllMap(adKey).Result()
		if ad == nil {
			continue
		}

		imp, _ := strconv.Atoi(ad["impressions"])
		data := &Report{
			&Ad{
				ad["slot"],
				ad["id"],
				ad["title"],
				ad["type"],
				ad["advertiser"],
				ad["destination"],
				imp,
			},
			0,
			imp,
			nil,
		}
		report[ad["id"]] = data
	}

	for adId, clicks := range getLog(advrId) {
		if _, exists := report[adId]; !exists {
			report[adId] = &Report{}
		}
		report[adId].Clicks = len(clicks)
	}
	c.JSON(200, report)
}

func routeGetFinalReport(c *gin.Context) {
	advrId := advertiserId(c.Request)

	if advrId == "" {
		c.Status(401)
		return
	}

	reports := map[string]*Report{}
	adKeys, _ := rd.SMembers(advertiserKey(advrId)).Result()
	for _, adKey := range adKeys {
		ad, _ := rd.HGetAllMap(adKey).Result()
		if ad == nil {
			continue
		}

		imp, _ := strconv.Atoi(ad["impressions"])
		data := &Report{
			&Ad{
				ad["slot"],
				ad["id"],
				ad["title"],
				ad["type"],
				ad["advertiser"],
				ad["destination"],
				imp,
			},
			0,
			imp,
			nil,
		}
		reports[ad["id"]] = data
	}

	logs := getLog(advrId)

	for adId, report := range reports {
		log, exists := logs[adId]
		if exists {
			report.Clicks = len(log)
		}

		breakdown := &BreakdownReport{
			map[string]int{},
			map[string]int{},
			map[string]int{},
		}
		for i := range log {
			click := log[i]
			incr_map(&breakdown.Gender, click.Gender)
			incr_map(&breakdown.Agents, click.Agent)
			generation := "unknown"
			if click.Age != -1 {
				generation = strconv.Itoa(click.Age / 10)
			}
			incr_map(&breakdown.Generations, generation)
		}
		report.Breakdown = breakdown
		reports[adId] = report
	}

	c.JSON(200, reports)
}
