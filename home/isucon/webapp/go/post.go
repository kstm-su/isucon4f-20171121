package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func routePostAd(r render.Render, req *http.Request, params martini.Params) {
	slot := params["slot"]

	advrId := advertiserId(req)
	if advrId == "" {
		r.Status(404)
		return
	}

	req.ParseMultipartForm(100000)
	asset := req.MultipartForm.File["asset"][0]
	id := nextAdId()
	key := adKey(slot, id)

	content_type := ""
	if len(req.Form["type"]) > 0 {
		content_type = req.Form["type"][0]
	}
	if content_type == "" && len(asset.Header["Content-Type"]) > 0 {
		content_type = asset.Header["Content-Type"][0]
	}
	if content_type == "" {
		content_type = "video/mp4"
	}

	title := ""
	if a := req.Form["title"]; a != nil {
		title = a[0]
	}
	destination := ""
	if a := req.Form["destination"]; a != nil {
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

	ioutil.WriteFile("/home/isucon/webapp/public/slots/"+slot+"/ads/"+id+"/asset.gz", gzipData, os.ModePerm)

	rd.RPush(slotKey(slot), id)
	rd.SAdd(advertiserKey(advrId), key)

	r.JSON(200, getAd(req, slot, id))
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
