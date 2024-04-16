package main

import (
	"math/rand"
	"time"

	"go-prometheus/internal/metrics2"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	router.Use(gin.Logger(), metrics2.Metrics())
	router.GET("/hello", func(context *gin.Context) {
		context.Set("result: ", "hello")
	})
	router.GET("/flight", func(context *gin.Context) {
		time.Sleep(1 * time.Second)
		context.Set("result: ", "success")
	})
	router.GET("/random", func(context *gin.Context) {
		randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10) + 1
		time.Sleep(time.Duration(randNum) * time.Second)
		context.Set("result: ", "success")
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(":8080")
}
