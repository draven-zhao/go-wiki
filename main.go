package main

import (
	"go-wiki/prometheus"
	"go-wiki/route"

	"github.com/gin-gonic/gin"
)

func main() {
	prometheus.StartScheduler()
	r := gin.Default()
	route.Api(r)
	r.Run(":8080")
}
