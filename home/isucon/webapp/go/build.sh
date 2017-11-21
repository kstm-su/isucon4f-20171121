#!/bin/bash

#go get github.com/go-martini/martini
go get github.com/gin-gonic/gin
#go get gopkg.in/go-redis/redis.v2
go get github.com/martini-contrib/render
go get github.com/kr/pretty
go get golang.org/x/sync/errgroup
go build -o golang-webapp .
