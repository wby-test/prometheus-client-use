package main

import (
	"time"

	"go-prometheus/internal/metrics2"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(gin.Logger(), metrics2.Metrics())
	router.GET("/hello", func(context *gin.Context) {
		context.Set("result: ", "hello")
	})
	router.GET("/flight", func(context *gin.Context) {
		time.Sleep(1 * time.Second)
		context.Set("result: ", "success")
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//metrics.AddHandleFunc(router)

	router.Run(":8088")
}
