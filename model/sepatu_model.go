package model

import (
	"github.com/google/uuid"
)

type Sepatu struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name          string    `json:"name" binding:"required"`
	Brand         string    `json:"brand" binding:"required"`
	Size          int       `json:"size" binding:"required"`
	Price         float64   `json:"price" binding:"required"`
	Stock         int       `json:"stock" binding:"required"`
	LastUpdatedAt int64     `json:"lastupdated_at"`
	CreatedAt     int64     `json:"created_at"`
}

type UpdateSepatu struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name          *string   `json:"name"`
	Brand         *string   `json:"brand"`
	Size          *int      `json:"size"`
	Price         *float64  `json:"price"`
	Stock         *int      `json:"stock"`
	LastUpdatedAt int64     `json:"lastupdated_at"`
	CreatedAt     int64     `json:"created_at"`
}

type DeleteSepatu struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id" binding:"required,uuid"`
}
