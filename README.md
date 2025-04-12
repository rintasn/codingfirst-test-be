# Sistem Manajemen Preferensi Pengguna

Sistem ini merupakan aplikasi manajemen preferensi pengguna yang memungkinkan pengguna untuk menyesuaikan pengalaman aplikasi mereka dengan mudah melalui antarmuka web atau perintah bahasa natural menggunakan Claude Desktop.

## Tech Stack

### Backend
- **Bahasa**: Golang
- **Framework/Library**:
  - `github.com/gorilla/mux` - Router HTTP
  - `github.com/rs/cors` - Middleware CORS
  - `github.com/xuri/excelize/v2` - Manipulasi Excel (untuk keperluan ekspor/impor)
  - `gorm.io/driver/postgres` - Driver PostgreSQL untuk GORM
  - `gorm.io/gorm` - ORM untuk Golang
  - `golang.org/x/crypto/bcrypt` - Enkripsi password
  - `github.com/golang-jwt/jwt/v4` - Implementasi JSON Web Token
  - `github.com/joho/godotenv` - Pembacaan file .env

### Frontend
- **Framework**: Next.js 15
- **Bahasa**: TypeScript
- **Styling**:
  - Tailwind CSS - Utility-first CSS framework
  - Shadcn UI - Komponen UI berbasis Tailwind
- **State Management**:
  - React Context API
- **Form & Validasi**:
  - React Hook Form
  - Zod

### Database
- PostgreSQL - Database relasional untuk penyimpanan data pengguna dan preferensi

## Fitur Utama

### 1. Autentikasi Pengguna
- Registrasi pengguna baru
- Login dengan JWT (JSON Web Token)
- Proteksi endpoint yang memerlukan autentikasi
- Enkripsi password dengan bcrypt

### 2. Manajemen Preferensi
- Pengaturan tema (terang/gelap)
- Pemilihan bahasa (Inggris, Spanyol, Indonesia)
- Pengaturan notifikasi (aktif/nonaktif)
- Penyimpanan preferensi di database

### 3. Context Management Protocol (MCP)
- Pemisahan data pengguna dan preferensi
- Penerapan preferensi secara global di aplikasi
- Perubahan real-time tanpa refresh halaman

### 4. Claude Desktop Integration
- Antarmuka chat untuk berinteraksi dengan AI
- Pemrosesan bahasa natural untuk mengubah preferensi
- Umpan balik instan saat preferensi diubah
- Saran perintah untuk pengguna baru

### 5. Antarmuka Pengguna yang Responsif
- Dashboard untuk melihat preferensi saat ini
- Halaman pengaturan untuk modifikasi manual
- Dukungan penuh untuk tema terang/gelap
- Desain responsif untuk berbagai ukuran perangkat

### 6. Persistensi Antar Sesi
- Preferensi disimpan di database
- Preferensi dimuat secara otomatis saat login
- Konsistensi pengalaman pengguna antar perangkat

## Alur Kerja Sistem

1. **Registrasi & Login**:
   - Pengguna mendaftar dengan username, email, dan password
   - Setelah login, JWT token dibuat dan preferensi default diterapkan

2. **Penggunaan Preferensi**:
   - Preferensi dimuat sebagai konteks dan diterapkan ke seluruh aplikasi
   - Tema, bahasa, dan pengaturan notifikasi diaplikasikan secara real-time

3. **Modifikasi Preferensi**:
   - Via halaman pengaturan dengan form UI
   - Via Claude Desktop dengan perintah bahasa natural
   - Perubahan disimpan ke database dan diterapkan segera

4. **Persistensi**:
   - Preferensi bertahan di berbagai sesi
   - Konsistensi pengalaman pengguna antar login

## Pengelolaan Database

### Skema Database dengan GORM

```go
// User Model
type User struct {
  ID        uint           `gorm:"primaryKey" json:"id"`
  Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
  Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
  Password  string         `gorm:"size:255;not null" json:"-"`
  Preferences UserPreferences `gorm:"foreignKey:UserID" json:"preferences"`
  CreatedAt time.Time      `json:"created_at"`
  UpdatedAt time.Time      `json:"updated_at"`
  DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserPreferences Model
type UserPreferences struct {
  ID        uint           `gorm:"primaryKey" json:"id"`
  UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
  Theme     string         `gorm:"size:20;default:'light'" json:"theme"`
  Language  string         `gorm:"size:20;default:'english'" json:"language"`
  Notifications bool        `gorm:"default:true" json:"notifications"`
  CreatedAt time.Time      `json:"created_at"`
  UpdatedAt time.Time      `json:"updated_at"`
}
```

## API Endpoints

### Autentikasi
- `POST /api/auth/register` - Pendaftaran pengguna baru
- `POST /api/auth/login` - Login pengguna

### Preferensi
- `GET /api/preferences` - Mengambil preferensi pengguna
- `POST /api/preferences` - Memperbarui preferensi pengguna

### Claude Desktop
- `POST /api/claude` - Mengirim pesan ke Claude dan menerima respons

### Pengguna
- `GET /api/user` - Mengambil data pengguna dengan preferensi

## Contoh Penggunaan Claude Desktop

Claude Desktop dapat memahami berbagai perintah bahasa natural, seperti:

- "Ubah tema ke mode gelap"
- "Ganti bahasa ke bahasa Indonesia"
- "Matikan notifikasi"
- "Apa pengaturan saya saat ini?"
- "Aktifkan notifikasi"
- "Kembali ke tema terang"

## Instalasi dan Penggunaan

### Backend
1. Clone repository
2. Install dependencies: `go mod tidy`
3. Buat file `.env` dengan konfigurasi database dan JWT
4. Jalankan aplikasi: `go run main.go`

### Frontend
1. Clone repository frontend
2. Install dependencies: `npm install`
3. Buat file `.env.local` dengan URL API backend
4. Jalankan aplikasi: `npm run dev`

## Kontribusi

Kontribusi dan saran untuk perbaikan sangat diterima. Silakan buat issue atau pull request untuk berkontribusi pada proyek ini.

## Lisensi

[MIT License](LICENSE)