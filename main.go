package main

import (
	"context"

	"github.com/ministryofjustice/cloud-platform-go-get-module/init_app"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/log"
)

func main() {
	ginMode, dataAddr, dataPassword, apiKey := init_app.InitEnvVars()
	dataClient := init_app.InitDataClient(dataAddr, dataPassword)

	g := new(errgroup.Group)

	g.Go(func() error {
		return init_app.InitData(dataClient)
	})

	if err := g.Wait(); err != nil {
		log.Errorf(context.Background(), "Error bootstraping the repo version data into redis: %v", err)
	}

	init_app.InitApi(dataClient, ginMode, apiKey)
}
