package main

import (
	"log"

	"github.com/VmesteApp/auth-service/config"
	"github.com/VmesteApp/auth-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't init config: %s", err)
	}

	app.Run(cfg)
}
