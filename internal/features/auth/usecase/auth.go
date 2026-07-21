package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/dto"
	userDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/domain"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/hash"
	jwtPkg "github.com/ramdhanrizkij/nest-glamping-api/pkg/jwt"
)

type usecase struct {
	authRepo      domain.Repository
	userRepo      userDomain.Repository
	jwtSecret     string
	jwtExpiry     time.Duration
	refreshSecret string
	refreshExpiry time.Duration
}

func NewUsecase(
	authRepo domain.Repository,
	userRepo userDomain.Repository,
	jwtSecret string,
	jwtExpiry time.Duration,
	refreshSecret string,
	refreshExpiry time.Duration,
) domain.Service {
	return &usecase{
		authRepo:      authRepo,
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiry:     jwtExpiry,
		refreshSecret: refreshSecret,
		refreshExpiry: refreshExpiry,
	}
}

func (u *usecase) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, _ := u.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, appErr.Conflict("email already registered")
	}

	role, err := u.authRepo.FindRoleByName("customer")
	if err != nil {
		return nil, appErr.Internal("default role not found")
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, appErr.Internal("failed to hash password")
	}

	user := &userDomain.User{
		ID:           uuid.New(),
		RoleID:       role.ID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		PhoneNumber:  req.Phone,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, appErr.Internal("failed to create user")
	}

	return u.generateTokens(user.ID.String(), role.ID.String(), user.Name, user.Email, user.PhoneNumber, role.Name)
}

func (u *usecase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, appErr.Unauthorized("invalid email or password")
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return nil, appErr.Unauthorized("invalid email or password")
	}

	role, err := u.authRepo.FindRoleByID(user.RoleID)
	if err != nil {
		return nil, appErr.Internal("role not found")
	}

	return u.generateTokens(user.ID.String(), role.ID.String(), user.Name, user.Email, user.PhoneNumber, role.Name)
}

func (u *usecase) RefreshToken(refreshToken string) (*dto.AuthResponse, error) {
	rt, err := u.authRepo.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, appErr.Unauthorized("invalid or expired refresh token")
	}

	claims, err := jwtPkg.ValidateToken(refreshToken, u.refreshSecret)
	if err != nil {
		_ = u.authRepo.DeleteRefreshToken(refreshToken)
		return nil, appErr.Unauthorized("invalid or expired refresh token")
	}

	user, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, appErr.Unauthorized("user not found")
	}

	role, err := u.authRepo.FindRoleByID(user.RoleID)
	if err != nil {
		return nil, appErr.Internal("role not found")
	}

	_ = u.authRepo.DeleteRefreshToken(rt.Token)

	return u.generateTokens(user.ID.String(), role.ID.String(), user.Name, user.Email, user.PhoneNumber, role.Name)
}

func (u *usecase) Logout(refreshToken string) error {
	return u.authRepo.DeleteRefreshToken(refreshToken)
}

func (u *usecase) generateTokens(userID, roleID, name, email, phone, roleName string) (*dto.AuthResponse, error) {
	accessToken, err := jwtPkg.GenerateToken(userID, roleID, u.jwtSecret, u.jwtExpiry)
	if err != nil {
		return nil, appErr.Internal("failed to generate access token")
	}

	refreshToken, err := jwtPkg.GenerateToken(userID, roleID, u.refreshSecret, u.refreshExpiry)
	if err != nil {
		return nil, appErr.Internal("failed to generate refresh token")
	}

	userUUID, _ := uuid.Parse(userID)
	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userUUID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(u.refreshExpiry),
	}
	if err := u.authRepo.SaveRefreshToken(rt); err != nil {
		return nil, appErr.Internal("failed to save refresh token")
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:          userID,
			Name:        name,
			Email:       email,
			PhoneNumber: phone,
			RoleID:      roleID,
			RoleName:    roleName,
		},
	}, nil
}
