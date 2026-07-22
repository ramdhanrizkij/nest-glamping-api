package main

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/ramdhanrizkij/nest-glamping-api/config"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/bootstrap"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/database"
)

// @title           Glamping API
// @version         1.0
// @description     REST API for glamping booking platform — manage tent types, amenities, bookings with dynamic pricing, and payments.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/ramdhanrizkij/nest-glamping-api/issues
// @contact.email  support@glamping.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"
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
