package startup

import (
	"os"

	"github.com/gammazero/workerpool"
	"github.com/pna/order-app-backend/internal"
	"github.com/pna/order-app-backend/internal/controller"
	"github.com/pna/order-app-backend/internal/database"
	log "github.com/sirupsen/logrus"
)

func Migrate() {
	// Open the database connection
	db := database.Open()

	database.MigrateUp(db)
}

func registerDependencies() *controller.ApiContainer {
	// Open database connection
	db := database.Open()

	return internal.InitializeContainer(db)
}

func Execute() {
	// Configure logrus
	log.SetOutput(os.Stdout)

	container := registerDependencies()

	wp := workerpool.New(2)

	wp.Submit(container.HttpServer.Run)

	wp.StopWait()
}
