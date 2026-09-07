package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetOne(r *gin.Engine, rdb utils.DataAccessLayer) {
	r.GET("/:repo", func(c *gin.Context) {
		repo := c.Param("repo")

		currentVersionFields, err := rdb.HGetAll(repo)
		if err != nil {
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Reading from Redis"},
			}
			utils.SendResponse(c, obj)
			return
		}

		if len(currentVersionFields) == 0 {
			obj := utils.Response{
				Status: http.StatusNotFound,
				Error:  []string{"Repo not found: " + repo},
			}
			utils.SendResponse(c, obj)
			return
		}

		obj := utils.Response{
			Status: http.StatusOK,
			Data: gin.H{
				"repo":           repo,
				"currentVersion": currentVersionFields["currentVersion"],
				"sha":            currentVersionFields["sha"],
			},
		}
		utils.SendResponse(c, obj)
	})
}
