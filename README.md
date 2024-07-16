# User Service

This project is a User Service written in Go using the Gin framework, MongoDB for data storage, and RabbitMQ for message publishing. The service provides both HTTP and gRPC APIs for user management, including creating, reading, updating, and deleting users. It also publishes messages to RabbitMQ when users are created, updated, or deleted so other service could subscribe to rabbitmq and listen for updates of users. User id will be sent as queue message content, then by id service can query user-service for user details.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Setup](#setup)
- [Configuration](#configuration)
- [Running the Service](#running-the-service)
- [HTTP API](#http-api)
- [gRPC API](#grpc-api)
- [Running Tests](#running-tests)

## Features

- Create, read, update, and delete users
- Pagination and filtering for listing users
- Publish messages to RabbitMQ on user creation, update, and deletion
- Support for both HTTP and gRPC APIs

## Requirements

- Docker
- Docker Compose
- Go 1.21 or later
- MongoDB
- RabbitMQ

## Setup

1. Build the Docker images and start the services:
    ```sh
    docker-compose up --build
    ```

## Configuration

The service is configured using environment variables. You can set these variables in the `.env` file or directly in the Docker Compose file. The key environment variables are:

- `MONGO_URI`: The URI for connecting to MongoDB.
- `DB`: The name of the MongoDB database.
- `RABBITMQ_URI`: The URI for connecting to RabbitMQ.

## Running the Service

To run the service, ensure that Docker and Docker Compose are installed and then use the following command:

```sh
docker-compose up
```

### The service will be available at:
```sh
HTTP: http://localhost:8080
gRPC: localhost:50051
```
### The service exposes the following HTTP endpoints:
```sh
POST /api/v1/users: Create a new user
GET /api/v1/users/:id: Get a user by ID
PUT /api/v1/users/:id: Update a user by ID
DELETE /api/v1/users/:id: Delete a user by ID
GET /api/v1/users: List users with pagination and filtering
Example HTTP request to create a user:
```

```sh
curl -X POST http://localhost:8080/api/v1/users -d '{
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "jdoe",
    "password": "secret",
    "email": "john@example.com",
    "country": "USA"
}' -H "Content-Type: application/json"
```

```sh
curl -X GET 'http://localhost:8080/api/v1/users?country=UK&page=1&page_size=10'
```

### gRPC API
#### The gRPC service exposes the following methods:
```sh
CreateUser(User) returns (UserID)
GetUser(UserID) returns (User)
UpdateUser(User) returns (Empty)
DeleteUser(UserID) returns (Empty)
ListUsers(ListUsersRequest) returns (ListUsersResponse)
```
The gRPC service is defined in the proto/user.proto file. To interact with the gRPC service, you can use a gRPC client like grpcurl.

Example gRPC request to create a user:

```sh
grpcurl -plaintext -d '{
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "jdoe",
    "password": "secret",
    "email": "john@example.com",
    "country": "USA"
}' localhost:50051 user.UserService/CreateUser
```

Running Tests
To run the tests for the service, use the following command:

```sh
go test ./...
```

### Made by: Marko Stojkovic (software.developer@hotmail.rs)