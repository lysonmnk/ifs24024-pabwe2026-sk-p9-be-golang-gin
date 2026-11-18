package auth

import (
	"errors"
	"strings"

	"app/internal/database"
	"app/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *RegisterRequest) (*UserProfileResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	GetProfile(userID uint) (*UserProfileResponse, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserProfileResponse, error)
}

type authService struct {
	repo      Repository
	secretKey string
}

func NewService(repo Repository, secretKey string) Service {
	return &authService{
		repo:      repo,
		secretKey: secretKey,
	}
}

func (s *authService) Register(req *RegisterRequest) (*UserProfileResponse, error) {
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	username := req.Username
	if username == "" {
		parts := strings.Split(req.Email, "@")
		username = parts[0]
	}

	existingUsername, _ := s.repo.FindByUsername(username)
	if existingUsername != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &database.User{
		Name:     req.Name,
		Username: username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "author",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user *database.User
	var err error

	if req.Email != "" {
		user, err = s.repo.FindByEmail(req.Email)
	}
	if user == nil && req.Username != "" {
		user, err = s.repo.FindByUsername(req.Username)
	}
	if user == nil && req.Email != "" && !strings.Contains(req.Email, "@") {
		user, err = s.repo.FindByUsername(req.Email)
	}
	if err != nil || user == nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	token, _ := middleware.GenerateToken(user.ID, user.Role, s.secretKey)

	return &AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		User: UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) GetProfile(userID uint) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	return &UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	if user.Username != req.Username {
		existingUser, _ := s.repo.FindByUsername(req.Username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username sudah digunakan")
		}
	}

	user.Name = req.Name
	user.Username = req.Username

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
