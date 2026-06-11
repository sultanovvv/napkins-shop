package main

import (
	"context"
	"log"

	app "napkins-shop/public-executor/init"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err = application.Start(ctx); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}