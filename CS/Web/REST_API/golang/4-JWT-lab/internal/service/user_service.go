package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/auth"
	"github.com/grainme/movie-api/internal/cache"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo repository.UserRepository
	rdb      *redis.Client
}

func NewUserService(repo repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		userRepo: repo,
		rdb:      rdb,
	}
}

func (s *UserService) Login(ctx context.Context, username, password string) (domain.UserResponse, error) {
	user, err := s.userRepo.FindUserByName(ctx, username)
	if err != nil {
		return domain.UserResponse{}, err
	}

	match, err := auth.ComparePassword(password, user.PasswordHash)
	if err != nil {
		return domain.UserResponse{}, err
	}

	if !match {
		return domain.UserResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return domain.UserResponse{}, err
	}

	refreshToken := auth.GenerateRefreshToken()
	// Caching the user infos needed to generate another access token
	err = cache.SetUserByRefreshToken(ctx, s.rdb, refreshToken, cache.UserCache{
		UserId:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		log.Printf("could not cache user: %v\n", err)
	}

	return domain.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken uuid.UUID) error {
	err := cache.DelUserByRefreshTokenId(ctx, s.rdb, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Register(ctx context.Context, userRequestArgs domain.CreateUserRequest) (domain.User, error) {
	user, err := s.userRepo.AddUser(ctx, userRequestArgs)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// the client sends this request: `POST /auth/refresh`
func (s *UserService) RefreshToken(ctx context.Context, refreshTokenId uuid.UUID) (domain.UserResponse, error) {
	user, err := cache.GetUserByRefreshToken(ctx, s.rdb, refreshTokenId)
	if err != nil {
		return domain.UserResponse{}, err
	}
	if user == nil {
		return domain.UserResponse{}, errors.New("refresh token not found")
	}

	accessToken, err := auth.GenerateAccessToken(user.UserId, user.Role)
	if err != nil {
		return domain.UserResponse{}, err
	}

	refreshToken := auth.GenerateRefreshToken()

	cache.SetUserByRefreshToken(ctx, s.rdb, refreshToken, cache.UserCache{
		UserId:   user.UserId,
		Username: user.Username,
		Role:     user.Role,
	})

	response := domain.UserResponse{
		ID:           user.UserId,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}

	err = cache.DelUserByRefreshTokenId(ctx, s.rdb, refreshTokenId)
	if err != nil {
		log.Printf("failed to delete user from cache: %v", err)
	}

	return response, nil
}
