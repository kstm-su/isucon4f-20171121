package main

import ()

func routeGetAdWithId(r render.Render, req *http.Request, params martini.Params) {
	slot := params["slot"]
	id := params["id"]
	ad := getAd(req, slot, id)
	if ad != nil {
		r.JSON(200, ad)
	} else {
		r.JSON(404, map[string]string{"error": "not_found"})
	}
}
