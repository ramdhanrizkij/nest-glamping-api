package main

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/ramdhanrizkij/nest-glamping-api/config"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/bootstrap"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/database"
)

func main() {
	config.LoadEnv()

	appCfg := config.NewAppConfig()
	dbCfg := config.NewDatabaseConfig()
	loggerCfg := config.NewLoggerConfig()
	loggerCfg.Setup()

	db, err := database.NewConnection(dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.RunMigration(db); err != nil {
		log.Fatal(err)
	}

	deps := bootstrap.NewDependencies(db)

	app := bootstrap.NewFiber()
	bootstrap.SetupRoutes(app, deps)

	log.Info(appCfg.Name, " starting on port ", appCfg.Port)
	if err := app.Listen(":" + appCfg.Port); err != nil {
		log.Fatal(err)
	}
}
