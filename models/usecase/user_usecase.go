package usecase

import (
	"backend/models"
	"backend/models/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"backend/controllers/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(id uint) error
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	CreateUser(req models.CreateUserRequest) (*models.User, error)
	UpdateUser(id uint, req *models.UpdateUserRequest, currentUserID uint, currentUserRole models.Role) (*models.User, error)
	DeleteUser(id uint, currentUserRole models.Role) error
	ChangePassword(userID uint, req *models.ChangePasswordRequest) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	if _, err := u.repo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := u.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     req.Role,
	}

	// Default role is user if not specified
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by username
	user, err := u.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *userUsecase) Logout(id uint) error {
	return nil
}

func (u *userUsecase) GetUserByID(id uint) (*models.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecase) GetAllUsers() ([]models.User, error) {
	return u.repo.GetAll()
}

func (u *userUsecase) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Check if username already exists
	if _, err := u.repo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := u.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone: req.Phone,
		Role:     req.Role,
	}

	// Default role is user if not specified
	if user.Role == "" {
		user.Role = models.RoleAdmin
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) UpdateUser(id uint, req *models.UpdateUserRequest, currentUserID uint, currentUserRole models.Role) (*models.User, error) {
	// Get user to update
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Authorization check
	// User can only update themselves, unless they're admin
	if currentUserRole != models.RoleAdmin && currentUserID != id {
		return nil, errors.New("unauthorized to update this user")
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	// Update fields
	if req.Email != "" {
		// Check if email already used by another user
		existingUser, err := u.repo.FindByEmail(req.Email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already used by another user")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}
	// Only admin can change role
	if req.Role != "" {
		if currentUserRole != models.RoleAdmin {
			return nil, errors.New("only admin can change user role")
		}
		user.Role = req.Role
	}

	if err := u.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) DeleteUser(id uint, currentUserRole models.Role) error {
	// Only admin can delete users
	if currentUserRole != models.RoleAdmin {
		return errors.New("only admin can delete users")
	}

	// Check if user exists
	_, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return u.repo.Delete(id)
}

func (u *userUsecase) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return u.repo.Update(user)
}

func generateToken(userID uint) (string, error) {
    // Token expires in 24 hours
    expirationTime := time.Now().Add(24 * time.Hour)
    
    // Create the JWT claims, which includes the user ID and expiry time
    claims := &middleware.CustomClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            // In JWT, the expiry time is expressed as unix milliseconds
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "your-app-name",
            Subject:   fmt.Sprintf("%d", userID),
        },
    }

    // Declare the token with the algorithm used for signing, and the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Get the secret key from environment variable
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET not configured")
    }
    
    // Create the JWT string
    return token.SignedString([]byte(secret))
}

