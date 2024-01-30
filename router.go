package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func HandleRequest(c *gin.Context) {
	minSecValue := int(55)
	if minSec, ok := os.LookupEnv("MIN_SEC"); ok {
		v, err := strconv.Atoi(minSec)
		if err != nil {
			fmt.Println("Can't convert this to an int!")
		}
		minSecValue = int(v)
	}
	maxSecValue := int(55)
	if maxSec, ok := os.LookupEnv("MAX_SEC"); ok {
		v, err := strconv.Atoi(maxSec)
		if err != nil {
			fmt.Println("Can't convert this to an int!")
		}
		maxSecValue = int(v)
	}
	sleepTime := randInt(minSecValue, maxSecValue)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	c.JSON(
		http.StatusOK,
		fmt.Sprintf("Hello from the example HTTP service. %v", sleepTime),
	)
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
