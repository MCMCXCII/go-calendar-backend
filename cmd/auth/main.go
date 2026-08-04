package main

import (
	"context"
	"log"
	"project/internal/auth/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
