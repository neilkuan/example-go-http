package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

// Simple example HTTP service for trying out Beyla.
// 20% of calls will fail with HTTP status 500.

func HandleRequest(c *gin.Context) {
	time.Sleep(time.Duration(rand.Float64()*400.0) * time.Millisecond)
	if rand.Int31n(100) < 80 {
		c.JSON(
			http.StatusOK,
			"Hello from the example HTTP service.",
		)
	} else {

		c.JSON(
			http.StatusInternalServerError,
			"Simulating an error response with HTTP status 500.",
		)
	}
}

func NewRedisConnect(host string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     string(host),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/random", HandleRequest)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")

		if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok {
			client := NewRedisConnect(redisHost)
			ctx := context.Background()
			val, err := client.Get(ctx, user).Result()
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"user": user, "value": val})
			defer client.Close()
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	})

	r.POST("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		var json struct {
			Value string `json:"value" binding:"required"`
		}
		if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok && c.Bind(&json) == nil {
			client := NewRedisConnect(redisHost)
			ctx := context.Background()
			_, err := client.Set(ctx, user, json.Value, 0).Result()
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusNotImplemented, gin.H{"user": user, "status": "save fail"})
				return
			}
			defer client.Close()
			c.JSON(http.StatusCreated, gin.H{"user": user, "value": json.Value})
			return
		}
		c.JSON(http.StatusNotImplemented, gin.H{"user": user, "status": "save fail"})
	})
	return r
}
