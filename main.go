package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ministryofjustice/cloud-platform-go-get-module/routes"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func initGin(rdbClient *redis.Client, apiKey string) *gin.Engine {
	r := gin.New()
	// TODO: switch to release mode for production
	// [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	// gin.SetMode(value string)
	routes.InitLogger(r)
	routes.InitRouter(r, rdbClient, apiKey)

	return r
}

func main() {
	redisVal, redisPresent := os.LookupEnv("REDIS_SECRET")
	if redisVal == "" || !redisPresent {
		log.Fatal("REDIS_SECRET is not set")
	}

	apiKeyVal, apiKeyPresent := os.LookupEnv("API_KEY")
	if apiKeyVal == "" || !apiKeyPresent {
		log.Fatal("API_KEY is not set")
	}

	rdbClient := utils.InitRedisClient(redisVal)
	r := initGin(rdbClient, apiKeyVal)

	// Listen and Server in 0.0.0.0:3000
	r.Run(":3000")
}
