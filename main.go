package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/handlers"
	"user-service/proto/user-service/proto"
	"user-service/repositories"
	"user-service/services"
)

func main() {

	_ = godotenv.Load(".env")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).
		SetServerAPIOptions(
			serverAPI)
	opts.SetMaxPoolSize(100)
	opts.SetMinPoolSize(10)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(os.Getenv("DB"))

	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	rabbitMQPublisher, err := services.NewRabbitMQPublisher(rabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQPublisher.Close()

	userRepository := repositories.NewUserRepository(db, "users")
	userService := services.NewUserService(userRepository, rabbitMQPublisher)

	userHandlerGrpc := handlers.NewUserHandlerGrpc(userService)

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userHandlerGrpc)
	reflection.Register(grpcServer)

	go func() {
		grpcListener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		defer grpcListener.Close()

		log.Println("gRPC server listening at :50051")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	httpRouter := gin.New()

	httpRouter.Use(gin.Logger())
	httpRouter.Use(gin.Recovery())

	userHandler := handlers.NewUserHandler(userService)

	v1 := httpRouter.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.ListUsers)
		v1.GET("/users/:id", userHandler.GetUserByID)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	go func() {
		httpAddr := ":8080"
		srv := &http.Server{
			Addr:    httpAddr,
			Handler: httpRouter,
		}

		log.Println("HTTP server listening at", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Servers gracefully stopped")
}
