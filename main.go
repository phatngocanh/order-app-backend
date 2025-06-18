package main

import (
	"os"

	_ "github.com/pna/order-app-backend/docs"

	"github.com/pna/order-app-backend/startup"
)

// @title Order App
// @version 1.0
// @description Order app
// @BasePath /api/v1
func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		startup.Migrate()
		return
	}

	startup.Execute()
}
