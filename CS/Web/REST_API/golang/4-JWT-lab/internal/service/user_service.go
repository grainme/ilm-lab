package service

import (
	"context"

	"github.com/grainme/movie-api/internal/auth"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/grainme/movie-api/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) Login(ctx context.Context, username, password string) (domain.User, error) {
	user, err := s.userRepo.FindUserByName(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	match, err := auth.ComparePassword(password, user.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}

	if !match {
		// the handler should check that user is empty!?
		return domain.User{}, nil
	}

	// we should return jwt accessToken?
	return user, nil
}

func (s *UserService) Register(ctx context.Context, userRequestArgs domain.CreateUserRequest) (domain.User, error) {
	user, err := s.userRepo.AddUser(ctx, userRequestArgs)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
