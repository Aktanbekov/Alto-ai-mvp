#!/bin/bash

echo "🧪 Running tests before starting servers..."

# Run Go tests from tests folder
echo "📝 Running Go tests..."
if ! go test ./tests/... -v; then
    echo "❌ Tests failed! Fix the issues before starting servers."
    exit 1
fi

echo "✅ All tests passed!"
echo ""
echo "🔄 Restarting backend and frontend..."

# Function to kill process on a port
kill_port() {
    local port=$1
    local pid=$(lsof -ti:$port)
    if [ ! -z "$pid" ]; then
        echo "🛑 Killing process on port $port (PID: $pid)..."
        kill -9 $pid 2>/dev/null || true
        sleep 1
    else
        echo "✅ Port $port is already free"
    fi
}

# Kill processes on ports 8080 (backend) and 5173 (frontend)
echo "🧹 Cleaning up ports..."
kill_port 8080
kill_port 5173

# Wait a moment for ports to be fully released
sleep 2

# Cleanup function to kill background processes on script exit
cleanup() {
    echo ""
    echo "🛑 Shutting down servers..."
    kill_port 8080
    kill_port 5173
    exit 0
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM EXIT

# Start backend
echo "🚀 Starting backend on port 8080..."
(go run cmd/api/main.go) &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "🚀 Starting frontend on port 5173..."
(cd frontend && npm run dev) &
FRONTEND_PID=$!

echo ""
echo "✅ Servers started!"
echo "   Backend:  http://localhost:8080"
echo "   Frontend: http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop all servers"

# Wait for background processes
wait
