package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetOne(r *gin.Engine, rdb utils.DataAccessLayer) {
	r.GET("/:repo", func(c *gin.Context) {
		repo := c.Param("repo")

		fields, err := rdb.HGetAll(repo)
		if err != nil {
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Reading from Redis"},
			}
			utils.SendResponse(c, obj)
			return
		}

		// HGETALL reports a missing key as an empty map rather than an error.
		if len(fields) == 0 {
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
				"currentVersion": fields["currentVersion"],
				"sha":            fields["sha"],
			},
		}
		utils.SendResponse(c, obj)
	})
}
