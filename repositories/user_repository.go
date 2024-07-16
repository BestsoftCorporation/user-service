package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"user-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	resp, err := repo.collection.InsertOne(ctx, user)
	user.ID = resp.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		log.Println("Error inserting user:", err)
	}
	return err
}

func (repo *UserRepository) ListUsers(ctx context.Context, filter bson.M, page int64, pageSize int64) ([]models.User, error) {
	var users []models.User

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := repo.collection.Find(ctx,
		filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		var raw bson.M
		err := cursor.Decode(&raw)
		if err != nil {
			return nil, err
		}
		err = cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		user.ID = raw["_id"].(primitive.ObjectID).Hex()
		user.CreatedAt = raw["createdat"].(primitive.DateTime).Time()
		user.UpdatedAt = raw["updatedat"].(primitive.DateTime).Time()
		users = append(users, user)
	}

	if len(users) == 0 {
		return []models.User{}, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id format", err)
		return nil, err
	}
	err = repo.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, id string, updateUser *models.User) error {
	update := bson.M{"$set": updateUser}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error updating user:", err)
	}

	_, err = repo.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Println("Error updating user:", err)
	}
	return err
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id format", err)
		return err
	}
	_, err = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Println("Error deleting user:", err)
	}
	return err
}
