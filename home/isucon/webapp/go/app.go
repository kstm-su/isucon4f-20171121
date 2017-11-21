package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v2"
)

type Ad struct {
	Slot        string `json:"slot"`
	Id          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Advertiser  string `json:"advertiser"`
	Destination string `json:"destination"`
	Impressions int    `json:"impressions"`
}

type AdWithEndpoints struct {
	Ad
	Asset    string `json:"asset"`
	Redirect string `json:"redirect"`
	Counter  string `json:"counter"`
}

type ClickLog struct {
	AdId   string `json:"ad_id"`
	User   string `json:"user"`
	Agent  string `json:"agent"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

type Report struct {
	Ad          *Ad              `json:"ad"`
	Clicks      int              `json:"clicks"`
	Impressions int              `json:"impressions"`
	Breakdown   *BreakdownReport `json:"breakdown,omitempty"`
}

type BreakdownReport struct {
	Gender      map[string]int `json:"gender"`
	Agents      map[string]int `json:"agents"`
	Generations map[string]int `json:"generations"`
}

var rd *redis.Client

func init() {
	rd = redis.NewTCPClient(&redis.Options{
		//Addr: "localhost:6379",
		Addr: "app1:6379",
		DB:   0,
	})
}

func getDir(name string) string {
	base_dir := "/tmp/go/"
	path := base_dir + name
	os.MkdirAll(path, 0755)
	return path
}

func urlFor(req *http.Request, path string) string {
	host := req.Host
	if host != "" {
		return "http://" + host + path
	} else {
		return path
	}
}

func fetch(hash map[string]string, key string, defaultValue string) string {
	if hash[key] == "" {
		return defaultValue
	} else {
		return hash[key]
	}
}

func incr_map(dict *map[string]int, key string) {
	_, exists := (*dict)[key]
	if !exists {
		(*dict)[key] = 0
	}
	(*dict)[key]++
}

func advertiserId(req *http.Request) string {
	return req.Header.Get("X-Advertiser-Id")
}

func adKey(slot string, id string) string {
	return "isu4:ad:" + slot + "-" + id
}

func assetKey(slot string, id string) string {
	return "isu4:asset:" + slot + "-" + id
}

func advertiserKey(id string) string {
	return "isu4:advertiser:" + id
}

func slotKey(slot string) string {
	return "isu4:slot:" + slot
}

func nextAdId() string {
	id, _ := rd.Incr("isu4:ad-next").Result()
	return strconv.FormatInt(id, 10)
}

func nextAd(req *http.Request, slot string) *AdWithEndpoints {
	key := slotKey(slot)
	id, _ := rd.RPopLPush(key, key).Result()
	if id == "" {
		return nil
	}
	ad := getAd(req, slot, id)
	if ad != nil {
		return ad
	} else {
		rd.LRem(key, 0, id).Result()
		return nil
		//		return nextAd(req, slot)
	}
}

func getAd(req *http.Request, slot string, id string) *AdWithEndpoints {
	key := adKey(slot, id)
	m, _ := rd.HGetAllMap(key).Result()

	if m == nil {
		return nil
	}
	if _, exists := m["id"]; !exists {
		return nil
	}

	imp, _ := strconv.Atoi(m["impressions"])
	path_base := "/slots/" + slot + "/ads/" + id
	var ad *AdWithEndpoints
	ad = &AdWithEndpoints{
		Ad{
			m["slot"],
			m["id"],
			m["title"],
			m["type"],
			m["advertiser"],
			m["destination"],
			imp,
		},
		urlFor(req, path_base+"/asset"),
		urlFor(req, path_base+"/redirect"),
		urlFor(req, path_base+"/count"),
	}
	return ad
}

func decodeUserKey(id string) (string, int) {
	if id == "" {
		return "unknown", -1
	}
	splitted := strings.Split(id, "/")
	gender := "male"
	if splitted[0] == "0" {
		gender = "female"
	}
	age, _ := strconv.Atoi(splitted[1])

	return gender, age
}

func routeGetAd(c *gin.Context) {
	slot := c.Param("slot")
	ad := nextAd(c.Request, slot)
	if ad != nil {
		c.Redirect(301, "/slots/"+slot+"/ads/"+ad.Id)
	} else {
		c.JSON(404, map[string]string{"error": "not_found"})
	}
}

func routeGetAdRedirect(c *gin.Context) {
	slot := c.Param("slot")
	id := c.Param("id")
	ad := getAd(c.Request, slot, id)

	if ad == nil {
		c.JSON(404, map[string]string{"error": "not_found"})
		return
	}

	isuad := ""
	cookie, err := c.Request.Cookie("isuad")
	if err != nil {
		if err != http.ErrNoCookie {
			panic(err)
		}
	} else {
		isuad = cookie.Value
	}
	ua := c.Request.Header.Get("User-Agent")

	path := getLogPath(ad.Advertiser)

	/*
		var f *os.File
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(f, "%s\t%s\t%s\n", ad.Id, isuad, ua)
		f.Close()
	*/
	rd.RPush(path, fmt.Sprintf("%s\t%s\t%s", ad.Id, isuad, ua))

	c.Redirect(301, ad.Destination)
}

func main() {
	r := gin.Default()

	{
		slots := r.Group("/slots/:slot")
		slots.POST("/ads", routePostAd)
		slots.GET("/ad", routeGetAd)
		slots.GET("/ads/:id", routeGetAdWithId)
		slots.POST("/ads/:id/count", routeGetAdCount)
		slots.GET("/ads/:id/redirect", routeGetAdRedirect)
	}

	{
		me := r.Group("/me")
		me.GET("/report", routeGetReport)
		me.GET("/final_report", routeGetFinalReport)
	}
	r.POST("/syncasset/:slot/:id", syncAsset)

	r.POST("/initialize", routePostInitialize)
	r.Run(":8080")
}
