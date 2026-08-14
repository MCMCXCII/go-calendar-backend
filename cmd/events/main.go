package main

import (
	"context"
	"log"
	"project/internal/events/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
