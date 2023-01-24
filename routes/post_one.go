package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

func InitPostOne(r *gin.Engine, rdb utils.DataAccessLayer, actualApiKey string) {
	r.POST("/update/:repo/:version", func(c *gin.Context) {
		repo := c.Param("repo")

		APIKey := c.Request.Header.Get("X-API-Key")

		if APIKey != actualApiKey {
			obj := utils.Response{
				Status: http.StatusUnauthorized,
				Error:  []string{"Invalid API Key supplied"},
			}

			utils.SendResponse(c, obj)
			return
		}

		if strings.TrimSpace(repo) == "" {
			obj := utils.Response{
				Status: http.StatusBadRequest,
				Error:  []string{"Repo parameter is must be provided `/update/:repo_name/:updated_version_number`"},
			}
			utils.SendResponse(c, obj)
			return
		}

		version := c.Param("version")
		if strings.TrimSpace(version) == "" {
			obj := utils.Response{
				Status: http.StatusBadRequest,
				Error:  []string{"Update parameter is must be provided eg. `/update/:repo_name/:updated_version_number`"},
			}
			utils.SendResponse(c, obj)
			return
		}

		err := rdb.Set(repo, version, 0).Err()
		// if there has been an error setting the value
		// handle the error
		if err != nil {
			fmt.Println(err)
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Writing to Redis"},
			}
			utils.SendResponse(c, obj)
			return
		}

		obj := utils.Response{
			Status:  http.StatusOK,
			Message: []string{repo + " updated to " + version},
		}
		utils.SendResponse(c, obj)
	})

}
