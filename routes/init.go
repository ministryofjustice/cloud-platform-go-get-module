package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
)

func InitRouter(r *gin.Engine, rdbClient *redis.Client, apiKey string) {
	InitGetAll(r, rdbClient)
	InitGetOne(r, rdbClient)
	InitPostOne(r, rdbClient, apiKey)
}

func InitLogger(r *gin.Engine) {
	logger, _ := zap.NewProduction()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

}
