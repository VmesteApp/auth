// @title           vmesteapp/auth-service
// @version         1.0
// @description     SSO service for vmesteapp.ru.
package main

import (
	"log"

	"github.com/VmesteApp/auth-service/config"
	_ "github.com/VmesteApp/auth-service/docs"
	"github.com/VmesteApp/auth-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't init config: %s", err)
	}

	app.Run(cfg)
}
