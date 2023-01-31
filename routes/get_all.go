package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitGetAll(r *gin.Engine, rdb utils.DataAccessLayer) {
	r.GET("/", func(c *gin.Context) {
		allRepoVersions, err := utils.GetAllRedisKeysAndValues(rdb)
		if err != nil {
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
