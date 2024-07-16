package handlers

import (
	"context"
	"user-service/models"
	"user-service/proto/user-service/proto"
	"user-service/services"

	"go.mongodb.org/mongo-driver/bson"
)

type UserHandlerGrpc struct {
	service *services.UserService
	proto.UnimplementedUserServiceServer
}

func NewUserHandlerGrpc(service *services.UserService) *UserHandlerGrpc {
	return &UserHandlerGrpc{service: service}
}

func (s *UserHandlerGrpc) mustEmbedUnimplementedUserServiceServer() {}

func (s *UserHandlerGrpc) CreateUser(ctx context.Context, req *proto.User) (*proto.UserID, error) {
	user := &models.User{

		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Password:  req.GetPassword(),
		Email:     req.GetEmail(),
		Country:   req.GetCountry(),
	}
	err := s.service.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &proto.UserID{Id: user.ID}, nil
}

func (s *UserHandlerGrpc) GetUser(ctx context.Context, req *proto.UserID) (*proto.User, error) {

	user, err := s.service.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (s *UserHandlerGrpc) UpdateUser(ctx context.Context, req *proto.User) (*proto.Empty, error) {

	updateUser := &models.User{
		ID:        req.GetId(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Password:  req.GetPassword(),
		Email:     req.GetEmail(),
		Country:   req.GetCountry(),
	}

	err := s.service.UpdateUser(ctx, req.GetId(), updateUser)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *UserHandlerGrpc) DeleteUser(ctx context.Context, req *proto.UserID) (*proto.Empty, error) {
	err := s.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *UserHandlerGrpc) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	filters := bson.M{}
	if req.Filter != nil {
		if req.Filter.FirstName != "" {
			filters["firstname"] = req.Filter.FirstName
		}
		if req.Filter.LastName != "" {
			filters["lastname"] = req.Filter.LastName
		}
		if req.Filter.Nickname != "" {
			filters["nickname"] = req.Filter.Nickname
		}
		if req.Filter.Country != "" {
			filters["country"] = req.Filter.Country
		}
	}

	users, err := s.service.ListUsers(ctx, filters, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.User
	for _, user := range users {
		protoUser := &proto.User{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
		}
		protoUsers = append(protoUsers, protoUser)
	}

	return &proto.ListUsersResponse{Users: protoUsers}, nil
}
