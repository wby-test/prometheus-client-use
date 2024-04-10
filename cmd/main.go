package main

import (
	"go-prometheus/internal/metrics"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(gin.Logger(), metrics.CountAndDuration())
	router.GET("/hello", func(context *gin.Context) {
		context.Set("result: ", "hello")
	})

	metrics.AddHandleFunc(router)

	router.Run(":8088")
}
