package init_app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ministryofjustice/cloud-platform-go-get-module/routes"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func initGin(rdbClient utils.DataAccessLayer, apiKey string, ginMode string) *gin.Engine {
	gin.SetMode(ginMode)

	r := gin.New()

	routes.InitLogger(r)

	routes.InitRouter(r, rdbClient, apiKey)

	return r
}

func initRedis(redisAddr, redisPassword string) utils.DataAccessLayer {
	redisOptions := &redis.Options{
		Addr:     redisAddr + ":6379",
		Password: redisPassword,
		DB:       0,
	}

	return utils.InitRedisClient(redisOptions)
}

func InitEnvVars() (string, string, string, string) {
	redisVal, redisPresent := os.LookupEnv("REDIS_SECRET")
	if redisVal == "" || !redisPresent {
		log.Fatal("REDIS_SECRET is not set")
	}

	apiKeyVal, apiKeyPresent := os.LookupEnv("API_KEY")
	if apiKeyVal == "" || !apiKeyPresent {
		log.Fatal("API_KEY is not set")
	}

	redisAddrVal, redisAddrPresent := os.LookupEnv("REDIS_ADDR")
	if redisAddrVal == "" || !redisAddrPresent {
		log.Fatal("REDIS_ADDR is not set")
	}

	ginMode := "debug"
	ginModeVal, ginModePresent := os.LookupEnv("GIN_MODE")
	if ginModeVal == "" || !ginModePresent {
		os.Setenv("GIN_MODE", ginMode)
		ginModeVal = ginMode
	}

	return ginModeVal, redisAddrVal, redisVal, apiKeyVal
}

func InitApi(dataClient utils.DataAccessLayer, ginMode, apiKey string) {
	r := initGin(dataClient, apiKey, ginMode)

	// Listen and Server in 0.0.0.0:3000
	err := r.Run(":3000")

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
