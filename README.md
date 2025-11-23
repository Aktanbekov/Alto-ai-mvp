# 🤖 AI Interviewer - Interview Practice Platform

> Master your interview skills with AI-powered practice sessions. Get personalized feedback, improve your answers, and ace your next interview.

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-19.1.1-61DAFB.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6.svg)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Features

- 🎯 **AI-Powered Practice** - Realistic interview questions tailored to your industry and role
- 💬 **Natural Conversation** - Experience realistic interview scenarios with context-aware AI
- 📈 **Progress Tracking** - Monitor your improvement with detailed analytics and insights
- ⚡ **Instant Feedback** - Get real-time feedback on your answers and suggestions for improvement
- 🔒 **Private & Secure** - Your interviews are completely confidential
- ⏰ **24/7 Availability** - Practice anytime, anywhere
- 🎨 **Modern UI** - Beautiful, responsive design with smooth animations

## 🚀 Tech Stack

### Backend
- **Go 1.24.2** - High-performance backend
- **Gin** - Web framework
- **PostgreSQL** - Database (optional, can use in-memory)
- **JWT** - Authentication
- **OAuth2** - Google authentication

### Frontend
- **React 19.1.1** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **React Router** - Routing

## 📋 Prerequisites

- **Go** 1.24.2 or higher
- **Node.js** 20.x or higher
- **npm** or **yarn**
- **PostgreSQL** (optional, for production)
- **Docker** (optional, for containerized deployment)

## 🛠️ Installation

### Option 1: Local Development

#### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/Aktanbekov/Alto-ai-mvp.git
cd altoai_mvp
```

2. Install Go dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Google OAuth
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# JWT
JWT_SECRET=your_jwt_secret_key

# Database (optional)
DATABASE_URL=postgres://user:password@localhost:5432/dbname

# Frontend URL (for production)
FRONTEND_URL=http://localhost:3000
```

4. Run the backend:
```bash
go run cmd/api/main.go
```

The backend will be available at `http://localhost:8080`

#### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

### Option 2: Docker Deployment (Recommended)

1. Build and run with Docker Compose:
```bash
docker-compose up --build
```

2. Access the application:
   - Frontend & Backend: `http://localhost:3000`

3. View logs:
```bash
docker-compose logs -f
```

4. Stop the containers:
```bash
docker-compose down
```

For more Docker details, see [DOCKER.md](DOCKER.md)

## 🎮 Usage

### Getting Started

1. **Visit the Landing Page**
   - Navigate to `http://localhost:5173` (dev) or `http://localhost:3000` (Docker)
   - Explore the features and benefits

2. **Sign In**
   - Click "Sign In" or "Get Started"
   - Use Google OAuth for quick authentication
   - Or use email/password (if implemented)

3. **Start Interview Practice**
   - After login, you'll be redirected to the Interview Practice chat
   - The AI interviewer will greet you and ask questions
   - Answer naturally and receive instant feedback

4. **Track Your Progress**
   - Monitor your interview progress with the progress bar
   - View message count and session time
   - Get tips and suggestions

## 📁 Project Structure

```
altoai_mvp/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── HomePage.jsx     # Landing page
│   │   │   ├── LoginPage.jsx    # Login page
│   │   │   └── Chat.tsx          # Interview practice chat
│   │   ├── components/
│   │   │   └── ProtectedRoute.tsx
│   │   ├── App.jsx               # Main app component
│   │   └── api.js                # API client
│   └── package.json
├── internal/
│   ├── auth/                     # Authentication handlers
│   ├── handlers/                 # HTTP handlers
│   ├── middleware/              # Middleware (CORS, JWT, etc.)
│   ├── models/                   # Data models
│   ├── repository/               # Data access layer
│   ├── router/                   # Route definitions
│   └── services/                 # Business logic
├── pkg/
│   ├── errors/                   # Error handling
│   └── response/                 # Response utilities
├── docker-compose.yml            # Docker Compose config
├── Dockerfile                    # Docker build file
├── go.mod                        # Go dependencies
└── README.md                     # This file
```

## 🔌 API Endpoints

### Authentication
- `GET /auth/google` - Initiate Google OAuth login
- `GET /auth/google/callback` - OAuth callback handler
- `GET /me` - Get current user info (protected)

### Health Check
- `GET /health` - Health check endpoint

### API v1
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## 🔐 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | Yes |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | Yes |
| `GOOGLE_REDIRECT_URL` | OAuth redirect URL | Yes |
| `JWT_SECRET` | Secret key for JWT tokens | Yes |
| `DATABASE_URL` | PostgreSQL connection string | No |
| `FRONTEND_URL` | Frontend URL for redirects | No |
| `GIN_MODE` | Gin mode (release/debug) | No |

## 🐳 Docker

The project includes full Docker support for easy deployment:

- **Multi-stage build** for optimized image size
- **Health checks** for container monitoring
- **Environment variable** support
- **Production-ready** configuration

See [DOCKER.md](DOCKER.md) for detailed Docker documentation.

## 🧪 Development

### Running Tests
```bash
# Backend tests
go test ./...

# Frontend tests (if configured)
cd frontend
npm test
```

### Building for Production

#### Backend
```bash
go build -o bin/api ./cmd/api
```

#### Frontend
```bash
cd frontend
npm run build
```

The production build will be in `frontend/dist/`

## 🚢 Deployment

### Using Docker (Recommended)

1. Build the image:
```bash
docker build -t altoai-mvp:latest .
```

2. Run the container:
```bash
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  --name altoai-mvp \
  altoai-mvp:latest
```

### Manual Deployment

1. Build the frontend:
```bash
cd frontend && npm run build
```

2. Build the backend:
```bash
go build -o bin/api ./cmd/api
```

3. Run the backend (it will serve the frontend):
```bash
./bin/api
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Your Name** - *Initial work*

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) and [React](https://reactjs.org/)
- UI design inspired by modern interview platforms
- Thanks to all contributors and users

## 📞 Support

For support, email support@altoai.com or open an issue in this repository.

---

⭐ If you find this project helpful, please give it a star!

