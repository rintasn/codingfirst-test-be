// handlers/claude_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"main/config"
	"main/models"
)

// ClaudeRequest merupakan struktur untuk permintaan ke Claude Desktop
type ClaudeRequest struct {
	Message string `json:"message"`
}

// ClaudeResponse merupakan struktur untuk respons dari Claude Desktop
type ClaudeResponse struct {
	Message     string                  `json:"message"`
	Preferences *models.UserPreferences `json:"preferences,omitempty"`
	Action      string                  `json:"action,omitempty"`
}

// ClaudeHandler menangani permintaan ke Claude Desktop
func ClaudeHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan ID pengguna dari konteks
	userID := r.Context().Value("userID").(uint)

	// Parse request body
	var req ClaudeRequest
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

	// Logika sederhana untuk memproses perintah bahasa natural
	message := strings.ToLower(req.Message)
	action := ""
	preferenceUpdated := false

	// Proses perintah untuk mengubah tema
	if strings.Contains(message, "theme") || strings.Contains(message, "tema") {
		if strings.Contains(message, "dark") || strings.Contains(message, "gelap") {
			preferences.Theme = "dark"
			preferenceUpdated = true
			action = "theme_updated"
		} else if strings.Contains(message, "light") || strings.Contains(message, "terang") {
			preferences.Theme = "light"
			preferenceUpdated = true
			action = "theme_updated"
		}
	}

	// Proses perintah untuk mengubah bahasa
	if strings.Contains(message, "language") || strings.Contains(message, "bahasa") {
		if strings.Contains(message, "english") || strings.Contains(message, "inggris") {
			preferences.Language = "english"
			preferenceUpdated = true
			action = "language_updated"
		} else if strings.Contains(message, "indonesia") || strings.Contains(message, "indonesian") {
			preferences.Language = "indonesia"
			preferenceUpdated = true
			action = "language_updated"
		} else if strings.Contains(message, "spanish") || strings.Contains(message, "spanyol") {
			preferences.Language = "spanish"
			preferenceUpdated = true
			action = "language_updated"
		}
	}

	// Proses perintah untuk mengubah notifikasi
	if strings.Contains(message, "notification") || strings.Contains(message, "notifikasi") {
		if strings.Contains(message, "on") || strings.Contains(message, "enable") || strings.Contains(message, "aktif") {
			preferences.Notifications = true
			preferenceUpdated = true
			action = "notifications_updated"
		} else if strings.Contains(message, "off") || strings.Contains(message, "disable") || strings.Contains(message, "nonaktif") {
			preferences.Notifications = false
			preferenceUpdated = true
			action = "notifications_updated"
		}
	}

	// Simpan perubahan jika ada yang diperbarui
	if preferenceUpdated {
		result = config.DB.Save(&preferences)
		if result.Error != nil {
			http.Error(w, "Failed to update preferences: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Buat pesan response yang sesuai
	responseMsg := ""
	if action == "theme_updated" {
		responseMsg = "I've updated your theme to " + preferences.Theme + " mode."
	} else if action == "language_updated" {
		responseMsg = "I've changed your language preference to " + preferences.Language + "."
	} else if action == "notifications_updated" {
		if preferences.Notifications {
			responseMsg = "I've turned notifications on for you."
		} else {
			responseMsg = "I've turned notifications off for you."
		}
	} else if strings.Contains(message, "what") && (strings.Contains(message, "theme") || strings.Contains(message, "tema")) {
		responseMsg = "Your current theme is set to " + preferences.Theme + " mode."
		action = "theme_info"
	} else if strings.Contains(message, "what") && (strings.Contains(message, "language") || strings.Contains(message, "bahasa")) {
		responseMsg = "Your current language is set to " + preferences.Language + "."
		action = "language_info"
	} else if strings.Contains(message, "notification") || strings.Contains(message, "notifikasi") && (strings.Contains(message, "status") || strings.Contains(message, "what")) {
		if preferences.Notifications {
			responseMsg = "Your notifications are currently enabled."
		} else {
			responseMsg = "Your notifications are currently disabled."
		}
		action = "notifications_info"
	} else {
		responseMsg = "I'm sorry, I don't understand that command. You can ask me to change your theme, language, or notification settings."
		action = "unknown_command"
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ClaudeResponse{
		Message:     responseMsg,
		Preferences: &preferences,
		Action:      action,
	})
}
