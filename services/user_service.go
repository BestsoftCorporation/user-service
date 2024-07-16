package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"user-service/models"
	"user-service/repositories"
)

type UserService struct {
	repo      *repositories.UserRepository
	publisher *RabbitMQPublisher
}

func NewUserService(repo *repositories.UserRepository, publisher *RabbitMQPublisher) *UserService {
	return &UserService{repo: repo, publisher: publisher}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := s.repo.CreateUser(ctx, user)
	if err == nil {
		s.publisher.Publish("user.created", user.ID)
	}
	return err
}

func (s *UserService) ListUsers(ctx context.Context, filters bson.M, page int64, pageSize int64) ([]models.User, error) {
	return s.repo.ListUsers(ctx, filters, page, pageSize)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, updateUser *models.User) error {
	now := time.Now()
	updateUser.UpdatedAt = now

	err := s.repo.UpdateUser(ctx, id, updateUser)
	if err == nil {
		s.publisher.Publish("user.updated", id)
	}
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.DeleteUser(ctx, id)
	if err == nil {
		s.publisher.Publish("user.deleted", id)
	}
	return err
}
