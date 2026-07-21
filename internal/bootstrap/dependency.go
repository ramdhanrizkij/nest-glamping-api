package bootstrap

import (
	"time"

	"github.com/ramdhanrizkij/nest-glamping-api/config"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities"
	amenitiesRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/repository"
	amenitiesUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/usecase"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth"
	authRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/repository"
	authUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/usecase"
	tenttypes "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types"
	tentTypesRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/repository"
	tentTypesUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/usecase"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents"
	tentsRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/repository"
	tentsUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/usecase"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users"
	usersRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/repository"
	usersUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/usecase"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB              *gorm.DB
	JWTSecret       string
	JWTExpiry       time.Duration
	RefreshSecret   string
	RefreshExpiry   time.Duration
	AuthModule      *auth.Module
	UserModule      *users.Module
	AmenityModule   *amenities.Module
	TentTypeModule  *tenttypes.Module
	TentModule      *tents.Module
}

func NewDependencies(db *gorm.DB) *Dependencies {
	jwtSecret := config.GetEnv("JWT_SECRET", "secret")
	jwtExpiry := config.GetEnv("JWT_EXPIRY", "15m")
	refreshSecret := config.GetEnv("JWT_REFRESH_SECRET", "refresh-secret")
	refreshExpiry := config.GetEnv("JWT_REFRESH_EXPIRY", "168h")

	expiry, _ := time.ParseDuration(jwtExpiry)
	refreshExp, _ := time.ParseDuration(refreshExpiry)

	// Repositories
	userRepository := usersRepo.NewRepository(db)
	authRepository := authRepo.NewRepository(db)
	amenityRepository := amenitiesRepo.NewRepository(db)
	tentTypeRepository := tentTypesRepo.NewRepository(db)
	tentRepository := tentsRepo.NewRepository(db)

	// Usecases
	userService := usersUsecase.NewUsecase(userRepository)
	authService := authUsecase.NewUsecase(authRepository, userRepository, jwtSecret, expiry, refreshSecret, refreshExp)
	amenityService := amenitiesUsecase.NewUsecase(amenityRepository)
	tentTypeService := tentTypesUsecase.NewUsecase(tentTypeRepository)
	tentService := tentsUsecase.NewUsecase(tentRepository, tentTypeRepository)

	// Modules
	userModule := users.NewModule(userService)
	authModule := auth.NewModule(authService)
	amenityModule := amenities.NewModule(amenityService)
	tentTypeModule := tenttypes.NewModule(tentTypeService)
	tentModule := tents.NewModule(tentService)

	return &Dependencies{
		DB:              db,
		JWTSecret:       jwtSecret,
		JWTExpiry:       expiry,
		RefreshSecret:   refreshSecret,
		RefreshExpiry:   refreshExp,
		AuthModule:      authModule,
		UserModule:      userModule,
		AmenityModule:   amenityModule,
		TentTypeModule:  tentTypeModule,
		TentModule:      tentModule,
	}
}
