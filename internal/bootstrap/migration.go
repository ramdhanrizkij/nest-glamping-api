package bootstrap

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v3/log"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/bootstrap/migrations"
	"gorm.io/gorm"
)

var allMigrations = []*gormigrate.Migration{
	migrations.CreateRoles,
	migrations.CreatePermissions,
	migrations.CreateUsers,
	migrations.CreateRefreshTokens,
	migrations.CreateTentTypes,
	migrations.CreateTentTypeRates,
	migrations.CreateTentTypeImages,
	migrations.CreateAmenities,
	migrations.CreateTentTypeAmenities,
	migrations.CreateTents,
	migrations.CreateBookings,
	migrations.CreateBookingTents,
	migrations.CreatePayments,
}

func RunMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, allMigrations)

	m.InitSchema(func(tx *gorm.DB) error {
		log.Info("Initializing schema (fresh database)...")
		if err := tx.AutoMigrate(
			&migrations.Role{},
			&migrations.Permission{},
			&migrations.RolePermission{},
			&migrations.User{},
			&migrations.RefreshToken{},
			&migrations.TentType{},
			&migrations.TentTypeRate{},
			&migrations.TentTypeImage{},
			&migrations.Amenity{},
			&migrations.TentTypeAmenity{},
			&migrations.Tent{},
			&migrations.Booking{},
			&migrations.BookingTent{},
			&migrations.Payment{},
		); err != nil {
			return err
		}

		roles := []migrations.Role{
			{Name: "admin", Description: "System administrator"},
			{Name: "customer", Description: "Regular customer"},
		}
		for _, r := range roles {
			tx.Where(migrations.Role{Name: r.Name}).FirstOrCreate(&r)
		}
		return nil
	})

	log.Info("Running migrations...")
	if err := m.Migrate(); err != nil {
		return err
	}
	log.Info("Migrations completed")
	return nil
}

func RollbackLast(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, allMigrations)
	return m.RollbackLast()
}

func RollbackAll(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, allMigrations)
	return m.RollbackTo("0")
}
