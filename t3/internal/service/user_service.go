package service

import (
	"context"
	"weather-api-t3/internal/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, name, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, name, email string) (*models.User, error)
	SoftDeleteUser(ctx context.Context, id int) error
	
	AddCity(ctx context.Context, userID int, cityName string) error
	ListCities(ctx context.Context, userID int) ([]string, error)
	DeleteCity(ctx context.Context, userID int, cityName string) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, name, email string) (*models.User, error) {
	return s.repo.CreateUser(ctx, name, email)
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *UserService) Get(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (*models.User, error) {
	return s.repo.UpdateUser(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.SoftDeleteUser(ctx, id)
}

func (s *UserService) AddCity(ctx context.Context, userID int, cityName string) error {
	return s.repo.AddCity(ctx, userID, cityName)
}

func (s *UserService) ListCities(ctx context.Context, userID int) ([]string, error) {
	return s.repo.ListCities(ctx, userID)
}

func (s *UserService) DeleteCity(ctx context.Context, userID int, cityName string) error {
	return s.repo.DeleteCity(ctx, userID, cityName)
}
