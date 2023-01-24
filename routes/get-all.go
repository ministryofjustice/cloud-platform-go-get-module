package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetAll(r *gin.Engine, rdb *redis.Client) {
	r.GET("/", func(c *gin.Context) {
		allRepoVersions, err := utils.GetAllRedisKeysAndValues(rdb)
		if err != nil {
			fmt.Println(err)
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Reading from Redis"},
			}
			utils.SendResponse(c, obj)
			return
		}

		obj := utils.Response{
			Status: http.StatusOK,
			Data:   allRepoVersions,
		}
		utils.SendResponse(c, obj)
	})

}
