package main

import (
)

func routeGetAdCount(r render.Render, params martini.Params) {
	slot := params["slot"]
	id := params["id"]
	key := adKey(slot, id)

	exists, _ := rd.Exists(key).Result()
	if !exists {
		r.JSON(404, map[string]string{"error": "not_found"})
		return
	}

	rd.HIncrBy(key, "impressions", 1).Result()
	r.Status(204)
}
