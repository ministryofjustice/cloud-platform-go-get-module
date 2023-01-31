package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetOne(r *gin.Engine, rdb utils.DataAccessLayer) {
	r.GET("/:repo", func(c *gin.Context) {
		repo := c.Param("repo")

		currentVersion, err := rdb.Get(repo)
		if err != nil {
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
