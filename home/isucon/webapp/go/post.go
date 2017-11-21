package main

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func routePostAd(c *gin.Context) {
	slot := c.Param("slot")

	advrId := advertiserId(c.Request)
	if advrId == "" {
		c.Status(404)
		return
	}

	c.Request.ParseMultipartForm(100000)
	asset := c.Request.MultipartForm.File["asset"][0]
	id := nextAdId()
	key := adKey(slot, id)

	content_type := ""
	if len(c.PostForm("type")) > 0 {
		content_type = c.Request.Form["type"][0]
	}
	if content_type == "" && len(asset.Header["Content-Type"]) > 0 {
		content_type = asset.Header["Content-Type"][0]
	}
	if content_type == "" {
		content_type = "video/mp4"
	}

	title := ""
	if a := c.Request.Form["title"]; a != nil {
		title = a[0]
	}
	destination := ""
	if a := c.Request.Form["destination"]; a != nil {
		destination = a[0]
	}
	rd.HMSet(key,
		"slot", slot,
		"id", id,
		"title", title,
		"type", content_type,
		"advertiser", advrId,
		"destination", destination,
		"impressions", "0",
	)

	f, _ := asset.Open()
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	//asset_data := string(buf.Bytes())

	// assetをファイルに書き出す
	//rd.Set(assetKey(slot, id), asset_data)
	if err := os.MkdirAll("/home/isucon/webapp/public/slots/"+slot+"/ads/"+id+"/", 0777); err != nil {
		return
	}

	gzipData, err := makeGzip(buf.Bytes())
	if err != nil {
		return
	}

	var g errgroup.Group

	g.Go(func() error {
		return ioutil.WriteFile("/home/isucon/webapp/public/slots/"+slot+"/ads/"+id+"/asset.gz", gzipData, os.ModePerm)
	})

	g.Go(func() error {
		gzipBuf := bytes.NewBuffer(gzipData)
		req, err := http.NewRequest("POST", "http://"+os.Getenv("OTHER1")+":8080/syncasset/"+slot+"/"+id, gzipBuf)
		if err != nil {
			return err
		}
		client := &http.Client{}
		_, err = client.Do(req)
		return err
	})

	g.Go(func() error {
		gzipBuf := bytes.NewBuffer(gzipData)
		req, err := http.NewRequest("POST", "http://"+os.Getenv("OTHER2")+":8080/syncasset/"+slot+"/"+id, gzipBuf)
		if err != nil {
			return err
		}
		client := &http.Client{}
		_, err = client.Do(req)
		return err
	})

	g.Go(func() error {
		err := rd.RPush(slotKey(slot), id).Err()
		if err != nil {
			return err
		}
		return rd.SAdd(advertiserKey(advrId), key).Err()
	})

	if err := g.Wait(); err != nil {
		log.Printf("Error! routePostAd() %#v\n", err)
	}

	c.JSON(200, getAd(c.Request, slot, id))
}

func syncAsset(c *gin.Context) {
	slot := c.Param("slot")
	id := c.Param("id")
	gzipData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.Request.Body.Close()

	if err := os.MkdirAll("/home/isucon/webapp/public/slots/"+slot+"/ads/"+id+"/", 0777); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	ioutil.WriteFile("/home/isucon/webapp/public/slots/"+slot+"/ads/"+id+"/asset.gz", gzipData, os.ModePerm)
}

func makeGzip(body []byte) ([]byte, error) {
	var b bytes.Buffer
	err := func() error {
		gw := gzip.NewWriter(&b)
		defer gw.Close()

		if _, err := gw.Write(body); err != nil {
			return err
		}
		return nil
	}()
	return b.Bytes(), err
}
