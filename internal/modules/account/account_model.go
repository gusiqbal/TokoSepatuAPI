package account

import "github.com/google/uuid"

// User mewakili tabel users di database.
// Sama sekali tidak ada tag `json` atau `binding` di sini.
type User struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	Name          string    `gorm:"type:varchar(255);not null"`
	UserName      string    `gorm:"type:varchar(100);unique;not null"`
	Email         string    `gorm:"type:varchar(100);unique;not null"`
	PasswordHash  string    `gorm:"type:varchar(255);not null"`
	PhoneNumber   string    `gorm:"type:varchar(20);not null"`
	Level         string    `gorm:"type:varchar(20);"`
	CreatedAt     int64     `gorm:"autoCreateTime:milli"`
	LastUpdatedAt int64     `gorm:"autoUpdateTime:milli"`
}

func GetUser() []any {
	return []any{
		&User{},
	}
}
