// handlers/preference_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"main/config"
	"main/models"
)

// UpdatePreferenceRequest merupakan struktur untuk memperbarui preferensi

type UpdatePreferenceRequest struct {
	Theme         *string `json:"theme,omitempty"`
	Language      *string `json:"language,omitempty"`
	Notifications *bool   `json:"notifications,omitempty"`
}

// GetPreferenceResponse merupakan struktur untuk respons preferensi
type GetPreferenceResponse struct {
	Preferences models.UserPreferences `json:"preferences"`
}

// GetPreferencesHandler menangani permintaan untuk mengambil preferensi pengguna
func GetPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan ID pengguna dari konteks (yang diset oleh AuthMiddleware)
	userID := r.Context().Value("userID").(uint)

	// Ambil preferensi dari database
	var preferences models.UserPreferences
	result := config.DB.Where("user_id = ?", userID).First(&preferences)
	if result.Error != nil {
		http.Error(w, "Failed to get preferences: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetPreferenceResponse{
		Preferences: preferences,
	})
}

// UpdatePreferencesHandler menangani permintaan untuk memperbarui preferensi pengguna
func UpdatePreferencesHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan ID pengguna dari konteks
	userID := r.Context().Value("userID").(uint)

	// Parse request body
	var req UpdatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ambil preferensi pengguna saat ini
	var preferences models.UserPreferences
	result := config.DB.Where("user_id = ?", userID).First(&preferences)
	if result.Error != nil {
		http.Error(w, "Failed to get preferences: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Update preferensi jika disediakan dalam request
	updated := false

	if req.Theme != nil {
		preferences.Theme = *req.Theme
		updated = true
	}

	if req.Language != nil {
		preferences.Language = *req.Language
		updated = true
	}

	if req.Notifications != nil {
		preferences.Notifications = *req.Notifications
		updated = true
	}

	// Simpan perubahan jika ada yang diperbarui
	if updated {
		result = config.DB.Save(&preferences)
		if result.Error != nil {
			http.Error(w, "Failed to update preferences: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetPreferenceResponse{
		Preferences: preferences,
	})
}

// GetUserHandler menangani permintaan untuk mengambil data pengguna dengan preferensi
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan ID pengguna dari konteks
	userID := r.Context().Value("userID").(uint)

	// Ambil data pengguna dengan preferensi
	var user models.User
	result := config.DB.Preload("Preferences").First(&user, userID)
	if result.Error != nil {
		http.Error(w, "Failed to get user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
