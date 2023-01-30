package main

import "github.com/ministryofjustice/cloud-platform-go-get-module/init_app"

func main() {
	ginMode, dataAddr, dataPassword, apiKey := init_app.InitEnvVars()
	dataClient := init_app.InitDataClient(dataAddr, dataPassword)
	go init_app.InitData(dataClient)
	init_app.InitApi(dataClient, ginMode, apiKey)
}
