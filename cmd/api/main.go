package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yz13-dev/checkhouse-remote-checker/internal/auth"
	api "github.com/yz13-dev/checkhouse-remote-checker/internal/checkers"
	"github.com/yz13-dev/checkhouse-remote-checker/internal/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	region_code := os.Getenv("REGION")

	router := gin.Default()

	{
		jwtService := auth.NewJWTService(
			os.Getenv("JWT_SECRET"),
		)
		v1 := router.Group("/v1")
		v1.Use(middleware.TokenAuthMiddleware(jwtService))
		v1.GET("/check/http", func(c *gin.Context) {

			url := c.Query("url")
			method := c.Query("method")
			timeout := c.Query("timeout")

			if url == "" || method == "" || timeout == "" {
				c.JSON(400, gin.H{
					"error": "url, method, and timeout are required",
				})
				return
			}

			metrics, err := api.CheckHttp(url)

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"metrics": metrics,
			})

		})
		v1.GET("/check/dns", func(c *gin.Context) {
			domain := c.Query("domain")

			if domain == "" {
				c.JSON(400, gin.H{
					"error": "domain is required",
				})
				return
			}

			metrics, err := api.CheckDNS(domain)

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"metrics": metrics,
			})
		})
		v1.GET("/check/tcp", func(c *gin.Context) {
			host := c.Query("host")
			port := c.Query("port")
			timeout := c.Query("timeout")

			if host == "" || port == "" || timeout == "" {
				c.JSON(400, gin.H{
					"error": "host, port, and timeout are required",
				})
				return
			}

			metrics, err := api.CheckTCP(host, port)

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"metrics": metrics,
			})
		})
	}

	router.GET("/health", func(c *gin.Context) {
		time := time.Now()

		c.JSON(200, gin.H{
			"timestamp":   time,
			"region_code": region_code,
		})
	})

	router.Run()
}
