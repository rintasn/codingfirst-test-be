# User Preferences Management System

This system is a user preferences management application that allows users to customize their application experience easily through a web interface or natural language commands using Claude Desktop.

## Tech Stack

### Backend
- **Language**: Golang
- **Framework/Libraries**:
  - `github.com/gorilla/mux` - HTTP Router
  - `github.com/rs/cors` - CORS Middleware
  - `github.com/xuri/excelize/v2` - Excel manipulation (for export/import purposes)
  - `gorm.io/driver/postgres` - PostgreSQL Driver for GORM
  - `gorm.io/gorm` - ORM for Golang
  - `golang.org/x/crypto/bcrypt` - Password encryption
  - `github.com/golang-jwt/jwt/v4` - JSON Web Token implementation
  - `github.com/joho/godotenv` - .env file reader

### Frontend
- **Framework**: Next.js 14
- **Language**: TypeScript
- **Styling**:
  - Tailwind CSS - Utility-first CSS framework
  - Shadcn UI - Tailwind-based UI components
- **State Management**:
  - React Context API
- **Form & Validation**:
  - React Hook Form
  - Zod

### Database
- PostgreSQL - Relational database for user data and preferences storage

## Main Features

### 1. User Authentication
- New user registration
- Login with JWT (JSON Web Token)
- Protection of endpoints requiring authentication
- Password encryption with bcrypt

### 2. Preferences Management
- Theme settings (light/dark)
- Language selection (English, Spanish, Indonesian)
- Notification settings (on/off)
- Preferences storage in database

### 3. Context Management Protocol (MCP)
- Separation of user data and preferences
- Global application of preferences
- Real-time changes without page refresh

### 4. Claude Desktop Integration
- Chat interface for AI interaction
- Natural language processing to change preferences
- Instant feedback when preferences are updated
- Command suggestions for new users

### 5. Responsive User Interface
- Dashboard to view current preferences
- Settings page for manual modification
- Full support for light/dark themes
- Responsive design for various device sizes

### 6. Cross-Session Persistence
- Preferences stored in database
- Preferences automatically loaded at login
- Consistent user experience across devices

## System Workflow

1. **Registration & Login**:
   - Users register with username, email, and password
   - After login, JWT token is created and default preferences are applied

2. **Preferences Usage**:
   - Preferences are loaded as context and applied throughout the application
   - Theme, language, and notification settings are applied in real-time

3. **Preferences Modification**:
   - Via settings page with UI form
   - Via Claude Desktop with natural language commands
   - Changes are saved to database and applied immediately

4. **Persistence**:
   - Preferences persist across various sessions
   - Consistent user experience between logins

## Database Management

### Database Schema with GORM

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

### Authentication
- `POST /api/auth/register` - New user registration
- `POST /api/auth/login` - User login

### Preferences
- `GET /api/preferences` - Retrieve user preferences
- `POST /api/preferences` - Update user preferences

### Claude Desktop
- `POST /api/claude` - Send message to Claude and receive response

### User
- `GET /api/user` - Retrieve user data with preferences

## Claude Desktop Usage Examples

Claude Desktop can understand various natural language commands, such as:

- "Change theme to dark mode"
- "Switch language to Indonesian"
- "Turn off notifications"
- "What are my current settings?"
- "Enable notifications"
- "Switch back to light theme"

## Installation and Usage

### Backend
1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Create `.env` file with database and JWT configuration
4. Run the application: `go run main.go`

### Frontend
1. Clone the frontend repository
2. Install dependencies: `npm install`
3. Create `.env.local` file with backend API URL
4. Run the application: `npm run dev`

## Contribution

Contributions and suggestions for improvements are welcome. Please create an issue or pull request to contribute to this project.

## License

[MIT License](LICENSE)