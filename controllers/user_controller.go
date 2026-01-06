package controllers

import (
	"backend/models"
	"backend/models/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) *UserController {
	return &UserController{usecase: usecase}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.usecase.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		Status:  http.StatusCreated,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.usecase.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *UserController) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Panggil usecase logout (opsional, jika perlu blacklist token)
	if err := ctrl.usecase.Logout(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "Logout successful",
	})
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := ctrl.usecase.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User profile retrieved successfully",
		Data:    user,
	})
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.usecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.usecase.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.usecase.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// Check if this is profile update or admin update
	isProfileUpdate, _ := c.Get("is_profile_update")
	
	var targetUserID uint
	
	if isProfileUpdate == true {
		// This is /users/profile route - updating own profile
		userID, _ := c.Get("user_id")
		targetUserID = userID.(uint)
	} else {
		// This is /admin/users/:id route
		idParam := c.Param("id")
		userID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		targetUserID = uint(userID)
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID, _ := c.Get("user_id")
	
	// Get user role - try from context first, fallback to user object
	var currentUserRole models.Role
	if role, exists := c.Get("user_role"); exists {
		currentUserRole = role.(models.Role)
	} else if user, exists := c.Get("user"); exists {
		currentUserRole = user.(models.User).Role
	} else {
		currentUserRole = models.RoleUser
	}

	user, err := ctrl.usecase.UpdateUser(targetUserID, &req, currentUserID.(uint), currentUserRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUser, _ := c.Get("user")
	currentUserRole := currentUser.(models.User).Role

	if err := ctrl.usecase.DeleteUser(uint(userID), currentUserRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	})
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	if err := ctrl.usecase.ChangePassword(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "Password changed successfully",
	})
}