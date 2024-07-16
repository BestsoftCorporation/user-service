package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
	"user-service/models"
	"user-service/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var updateUser models.User
	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdateUser(c.Request.Context(), userID, &updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	filters := bson.M{}

	if firstName := c.Query("first_name"); firstName != "" {
		filters["firstname"] = firstName
	}
	if lastName := c.Query("last_name"); lastName != "" {
		filters["lastname"] = lastName
	}
	if nickname := c.Query("nickname"); nickname != "" {
		filters["nickname"] = nickname
	}
	if country := c.Query("country"); country != "" {
		filters["country"] = country
	}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.service.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
