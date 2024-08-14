package main

import (
	"fmt"
	"log"

	"github.com/VmesteApp/auth-service/config"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("can't init config: %s", err)
	}

	fmt.Println(cfg)
}
