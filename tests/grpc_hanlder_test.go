package tests

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
	"user-service/proto/user-service/proto"
)

func TestCreateUser(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user := &proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "jdoe",
		Password:  "secret",
		Email:     "john@example.com",
		Country:   "USA",
	}

	resp, err := client.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	log.Printf("Created user with ID: %s", resp.GetId())
}

func TestUpdateUser(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	user := &proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "jdoe",
		Password:  "secret",
		Email:     "john@example.com",
		Country:   "USA",
	}

	resp, err := client.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	print(resp.GetId())
	updateUser := &proto.User{
		Id:        resp.GetId(),
		FirstName: "John Updated",
		LastName:  "Doe Updated",
	}

	_, err = client.UpdateUser(ctx, updateUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	getResp, err := client.GetUser(ctx, &proto.UserID{Id: resp.GetId()})
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if getResp.GetFirstName() != updateUser.GetFirstName() ||
		getResp.GetLastName() != updateUser.GetLastName() ||
		getResp.GetNickname() != updateUser.GetNickname() ||
		getResp.GetPassword() != updateUser.GetPassword() ||
		getResp.GetEmail() != updateUser.GetEmail() ||
		getResp.GetCountry() != updateUser.GetCountry() {
		t.Fatalf("User details do not match updated details")
	}

	log.Printf("Successfully updated user with ID: %s", resp.GetId())
}

func TestListUsers(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)

	// Create multiple users
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	users := []*proto.User{
		{
			FirstName: "Alice",
			LastName:  "Smith",
			Nickname:  "asmith",
			Password:  "password1",
			Email:     "alice@example.com",
			Country:   "USA",
		},
		{
			FirstName: "Bob",
			LastName:  "Jones",
			Nickname:  "bjones",
			Password:  "password2",
			Email:     "bob@example.com",
			Country:   "Canada",
		},
	}

	for _, user := range users {
		_, err := client.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// List users with pagination
	listResp, err := client.ListUsers(ctx, &proto.ListUsersRequest{})
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(listResp.GetUsers()) < len(users) {
		t.Fatalf("Expected at least %d users, got %d", len(users), len(listResp.GetUsers()))
	}

	for _, user := range listResp.GetUsers() {
		log.Printf("Listed user: %v", user)
	}

	log.Printf("Successfully listed users")
}

func TestDeleteUser(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	user := &proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "jdoe",
		Password:  "secret",
		Email:     "john@example.com",
		Country:   "USA",
	}

	createResp, err := client.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	_, err = client.DeleteUser(ctx, &proto.UserID{Id: createResp.GetId()})
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = client.GetUser(ctx, &proto.UserID{Id: createResp.GetId()})
	if err == nil {
		t.Fatalf("User was not deleted")
	}

	log.Printf("Successfully deleted user with ID: %s", createResp.GetId())
}
