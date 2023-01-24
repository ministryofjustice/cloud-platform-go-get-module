package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetOne(r *gin.Engine, rdb *redis.Client) {
	r.GET("/:repo", func(c *gin.Context) {
		repo := c.Param("repo")

		currentVersion, err := rdb.Get(repo).Result()
		if err != nil {
			fmt.Println(err)
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Reading from Redis, check the repo_name param for typos: " + repo},
			}
			utils.SendResponse(c, obj)
			return
		}

		obj := utils.Response{
			Status: http.StatusOK,
			Data: gin.H{
				"repo":           repo,
				"currentVersion": currentVersion,
			},
		}
		utils.SendResponse(c, obj)
	})

}
