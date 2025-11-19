package services

import (
	"context"
	"errors"
	"time"

	"altoai_mvp/internal/models"
	"altoai_mvp/internal/repository"
)

type UserService interface {
	List(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context, id string) (models.User, error)
	Create(ctx context.Context, dto models.CreateUserDTO) (models.User, error)
	Update(ctx context.Context, id string, dto models.UpdateUserDTO) (models.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}

func (s *userService) List(ctx context.Context) ([]models.User, error) {
	ctx, cancel := s.withTimeout(ctx); defer cancel()
	return s.repo.List()
}

func (s *userService) Get(ctx context.Context, id string) (models.User, error) {
	ctx, cancel := s.withTimeout(ctx); defer cancel()
	return s.repo.Get(id)
}

func (s *userService) Create(ctx context.Context, dto models.CreateUserDTO) (models.User, error) {
	ctx, cancel := s.withTimeout(ctx); defer cancel()
	return s.repo.Create(dto.Email, dto.Name)
}

func (s *userService) Update(ctx context.Context, id string, dto models.UpdateUserDTO) (models.User, error) {
	ctx, cancel := s.withTimeout(ctx); defer cancel()
	return s.repo.Update(id, dto.Email, dto.Name)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	ctx, cancel := s.withTimeout(ctx); defer cancel()
	return s.repo.Delete(id)
}

var ErrNotFound = errors.New("not found") // you can map repo errors if needed