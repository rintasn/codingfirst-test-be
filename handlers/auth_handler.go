// handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"main/config"
	"main/models"
	"main/utils"
)

// RegisterRequest merupakan struktur untuk permintaan registrasi
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest merupakan struktur untuk permintaan login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse merupakan struktur untuk respons autentikasi
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// RegisterHandler menangani permintaan registrasi pengguna baru
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Cek apakah username atau email sudah digunakan
	var existingUser models.User
	result := config.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	// Buat user baru
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simpan user ke database
	result = config.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Failed to create user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Buat preferensi default untuk user
	preferences := models.UserPreferences{
		UserID:        user.ID,
		Theme:         "light",
		Language:      "english",
		Notifications: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Simpan preferensi ke database
	result = config.DB.Create(&preferences)
	if result.Error != nil {
		http.Error(w, "Failed to create user preferences: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Load preferensi untuk response
	config.DB.Model(&user).Association("Preferences").Find(&user.Preferences)

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

// LoginHandler menangani permintaan login pengguna
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Cari user di database
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.RowsAffected == 0 {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Cek password
	if !user.CheckPassword(req.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Load preferensi untuk response
	config.DB.Preload("Preferences").First(&user, user.ID)

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}
