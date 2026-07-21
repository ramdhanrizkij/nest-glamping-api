package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/ramdhanrizkij/nest-glamping-api/config"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/bootstrap"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/database"
)

func main() {
	action := flag.String("action", "migrate", "migrate | rollback | rollback-all | seed")
	flag.Parse()

	config.LoadEnv()

	dbCfg := config.NewDatabaseConfig()
	db, err := database.NewConnection(dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "migrate":
		if err := bootstrap.RunMigration(db); err != nil {
			log.Fatal("Migration failed: ", err)
		}
	case "seed":
		if err := bootstrap.Seed(db); err != nil {
			log.Fatal("Seed failed: ", err)
		}
	case "rollback":
		if err := bootstrap.RollbackLast(db); err != nil {
			log.Fatal("Rollback failed: ", err)
		}
		log.Info("Rollback last migration completed")
	case "rollback-all":
		if err := bootstrap.RollbackAll(db); err != nil {
			log.Fatal("Rollback all failed: ", err)
		}
		log.Info("Rollback all completed")
	default:
		fmt.Println("Usage: go run cmd/migrate/main.go -action=[migrate|rollback|rollback-all|seed]")
		os.Exit(1)
	}
}
