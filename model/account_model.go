package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name          string    `json:"name" binding:"required"`
	UserName      string    `json:"userName" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	PasswordHash  string    `json:"passwordHash" binding:"required"`
	PhoneNumber   string    `json:"phoneNumber" binding:"required"`
	CreatedAt     int64     `json:"createdAt"`
	LastUpdatedAt int64     `json:"lastUpdatedAt"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
