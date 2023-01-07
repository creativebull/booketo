package users

import (
	"time"

	"github.com/google/uuid"
)

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID     `json:"session_id"`
	AccessToken           string        `json:"access_token"`
	AccessTokenExpiresAt  time.Time     `json:"access_token_expires_at"`
	RefreshToken          string        `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time     `json:"refresh_token_expires_at"`
	User                  *UserResponse `json:"user"`
}
