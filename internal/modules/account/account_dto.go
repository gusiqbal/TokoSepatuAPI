package account

type RegisterUserRequest struct {
	Name        string `json:"name" binding:"required"`
	UserName    string `json:"userName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phoneNumber"`
	Level       string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	CreatedAt   int64  `json:"createdAt"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateProfileRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"omitempty,email"`
	PhoneNumber *string `json:"phoneNumber"`
}
