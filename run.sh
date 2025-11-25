#!/bin/bash

echo "Restarting backend and frontend..."

# start backend
(go run cmd/api/main.go) &

# start frontend
(cd frontend && npm run dev) &
