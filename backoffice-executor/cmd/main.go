package main

import (
	"context"
	"log"

	app "napkins-shop/backoffice-executor/init"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Start(context.Background()); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
