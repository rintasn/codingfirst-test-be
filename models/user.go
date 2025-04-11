// Struktur Proyek:
//
// project/
// ├── main.go
// ├── config/
// │   └── database.go
// ├── models/
// │   └── user.go
// ├── handlers/
// │   └── auth_handler.go
// ├── middleware/
// │   └── auth_middleware.go
// └── utils/
//     └── jwt.go

// models/user.go
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User merupakan model untuk tabel users di database
type User struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Username    string          `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email       string          `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password    string          `gorm:"size:255;not null" json:"-"` // Password tidak ditampilkan dalam JSON response
	Preferences UserPreferences `gorm:"foreignKey:UserID" json:"preferences"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

// UserPreferences menyimpan preferensi pengguna
type UserPreferences struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Theme         string    `gorm:"size:20;default:'light'" json:"theme"`      // light/dark
	Language      string    `gorm:"size:20;default:'english'" json:"language"` // english/spanish/indonesia/etc
	Notifications bool      `gorm:"default:true" json:"notifications"`         // true/false
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel untuk model User
func (User) TableName() string {
	return "users"
}

// TableName menentukan nama tabel untuk model UserPreferences
func (UserPreferences) TableName() string {
	return "user_preferences"
}

// BeforeSave - hook yang dijalankan sebelum menyimpan user untuk hash password
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Hash password jika password diubah dan tidak kosong
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword membandingkan password yang diberikan dengan password hash yang tersimpan
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
