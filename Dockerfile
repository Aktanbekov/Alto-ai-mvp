# Multi-stage Dockerfile for AltoAI MVP
# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install frontend dependencies
RUN npm ci

# Copy frontend source code
COPY frontend/ ./

# Build frontend for production
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 3: Final Runtime Image
FROM alpine:latest

# Install ca-certificates and wget for HTTPS requests and health checks
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Copy the built Go binary from backend-builder
COPY --from=backend-builder /app/main .

# Copy the built frontend from frontend-builder
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy .env file if it exists (optional, can be overridden with env vars)
# COPY .env .env

# Expose port 8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]

