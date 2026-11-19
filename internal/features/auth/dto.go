package auth

import "time"

type RegisterRequest struct {
	Name     string `json:"name" binding:"max=100"`
	Username string `json:"username" binding:"max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Identifier string `json:"identifier"`
	Password   string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"max=100"`
	Username string `json:"username" binding:"max=50"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AuthResponse struct {
	Token     string              `json:"token"`
	TokenType string              `json:"token_type"`
	User      UserProfileResponse `json:"user"`
}
