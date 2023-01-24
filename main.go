package main

import "github.com/ministryofjustice/cloud-platform-go-get-module/init_app"

func main() {
	ginMode, redisAddr, redisPassword, apiKey := init_app.InitEnvVars()
	init_app.InitApi(ginMode, redisAddr, redisPassword, apiKey)
}
