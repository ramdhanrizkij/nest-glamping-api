package bootstrap

import (
	"time"

	"github.com/ramdhanrizkij/nest-glamping-api/config"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth"
	authRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/repository"
	authUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/usecase"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users"
	usersRepo "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/repository"
	usersUsecase "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/usecase"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB            *gorm.DB
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshSecret string
	RefreshExpiry time.Duration
	AuthModule    *auth.Module
	UserModule    *users.Module
}

func NewDependencies(db *gorm.DB) *Dependencies {
	jwtSecret := config.GetEnv("JWT_SECRET", "secret")
	jwtExpiry := config.GetEnv("JWT_EXPIRY", "15m")
	refreshSecret := config.GetEnv("JWT_REFRESH_SECRET", "refresh-secret")
	refreshExpiry := config.GetEnv("JWT_REFRESH_EXPIRY", "168h")

	expiry, _ := time.ParseDuration(jwtExpiry)
	refreshExp, _ := time.ParseDuration(refreshExpiry)

	userRepo := usersRepo.NewRepository(db)
	userService := usersUsecase.NewUsecase(userRepo)

	authRepository := authRepo.NewRepository(db)
	authService := authUsecase.NewUsecase(authRepository, userRepo, jwtSecret, expiry, refreshSecret, refreshExp)

	userModule := users.NewModule(userService)
	authModule := auth.NewModule(authService)

	return &Dependencies{
		DB:            db,
		JWTSecret:     jwtSecret,
		JWTExpiry:     expiry,
		RefreshSecret: refreshSecret,
		RefreshExpiry: refreshExp,
		AuthModule:    authModule,
		UserModule:    userModule,
	}
}
