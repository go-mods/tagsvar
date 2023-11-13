package _example

import "time"

type User struct {
	// User information
	ID        int    `json:"id"         gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string `json:"email"      gorm:"uniqueIndex;not null"`
	Name      string `json:"name"       gorm:"not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name"  gorm:"not null"`
	Password  string `json:"-"          gorm:"not null"`
	Role      Role   `json:"role"       gorm:"type:varchar(255);not null"`

	// Timestamps
	CreatedAt time.Time  `json:"-"      gorm:"not null"`
	UpdatedAt time.Time  `json:"-"      gorm:"not null"`
	DeletedAt *time.Time `json:"-"      gorm:"index"`
}
