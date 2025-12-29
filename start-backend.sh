#!/bin/bash

# Load environment variables from .env.local if it exists
if [ -f .env.local ]; then
    echo "📝 Loading environment variables from .env.local..."
    export $(cat .env.local | grep -v '^#' | xargs)
fi

# Set defaults for local development
export GIN_MODE=debug
export FRONTEND_URL=${FRONTEND_URL:-http://localhost:5173}
export GOOGLE_REDIRECT_URL=${GOOGLE_REDIRECT_URL:-http://localhost:8080/auth/google/callback}

# PostgreSQL connection (adjust if needed)
export POSTGRES_HOST=${POSTGRES_HOST:-localhost}
export POSTGRES_PORT=${POSTGRES_PORT:-5432}
export POSTGRES_USER=${POSTGRES_USER:-$(whoami)}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}
export POSTGRES_DB=${POSTGRES_DB:-altoai_db}

echo "🚀 Starting backend on port 8080..."
echo "   FRONTEND_URL: $FRONTEND_URL"
echo "   GOOGLE_REDIRECT_URL: $GOOGLE_REDIRECT_URL"
echo "   POSTGRES_HOST: $POSTGRES_HOST:$POSTGRES_PORT"

go run cmd/api/main.go
