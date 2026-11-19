package auth

import (
	"errors"
	"net/mail"
	"strings"

	"app/internal/database"
	"app/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
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

func (s *authService) Register(req *RegisterRequest) (*AuthResponse, error) {
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

	name := req.Name
	if name == "" {
		name = username
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &database.User{
		Name:     name,
		Username: username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "member",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
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

func (s *authService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user *database.User

	var candidates []string
	if req.Email != "" {
		candidates = append(candidates, req.Email)
	}
	if req.Username != "" && req.Username != req.Email {
		candidates = append(candidates, req.Username)
	}
	if req.Identifier != "" && req.Identifier != req.Email && req.Identifier != req.Username {
		candidates = append(candidates, req.Identifier)
	}

	for _, cand := range candidates {
		cand = strings.TrimSpace(cand)
		if cand == "" {
			continue
		}
		if strings.Contains(cand, "@") {
			user, _ = s.repo.FindByEmail(cand)
			if user == nil {
				user, _ = s.repo.FindByUsername(cand)
			}
		} else {
			user, _ = s.repo.FindByUsername(cand)
			if user == nil {
				user, _ = s.repo.FindByEmail(cand)
			}
		}
		if user != nil {
			break
		}
	}

	if user == nil {
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

	if req.Username != "" && user.Username != req.Username {
		existingUser, _ := s.repo.FindByUsername(req.Username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username sudah digunakan")
		}
		user.Username = req.Username
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		addr, err := mail.ParseAddress(req.Email)
		if err != nil || !strings.Contains(addr.Address, ".") {
			return nil, errors.New("format email tidak valid")
		}
		if user.Email != req.Email {
			existingUser, _ := s.repo.FindByEmail(req.Email)
			if existingUser != nil && existingUser.ID != userID {
				return nil, errors.New("email sudah terdaftar")
			}
			user.Email = req.Email
		}
	}

	if req.Password != "" {
		if len(req.Password) < 8 {
			return nil, errors.New("password minimal 8 karakter")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

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
