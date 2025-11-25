# How to Run AltoAI MVP with Docker

## Quick Start (Recommended - Using Docker Compose)

### 1. Build and run the application:
```bash
docker-compose up --build
```

This will:
- Build the Docker image
- Start the container
- Show logs in the terminal

### 2. Access the application:
Open your browser and go to: **http://localhost:8080**

### 3. Stop the application:
Press `Ctrl+C` in the terminal, or run:
```bash
docker-compose down
```

### 4. Run in background (detached mode):
```bash
docker-compose up -d --build
```

### 5. View logs:
```bash
docker-compose logs -f
```

### 6. Stop background container:
```bash
docker-compose down
```

---

## Alternative: Using Docker Directly

### 1. Build the Docker image:
```bash
docker build -t altoai-mvp:latest .
```

### 2. Run the container:
```bash
docker run -d \
  --name altoai-mvp \
  -p 8080:8080 \
  --env-file .env \
  --restart unless-stopped \
  altoai-mvp:latest
```

### 3. View logs:
```bash
docker logs -f altoai-mvp
```

### 4. Stop the container:
```bash
docker stop altoai-mvp
```

### 5. Remove the container:
```bash
docker rm altoai-mvp
```

---

## Verify It's Running

### Check health endpoint:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"ok":true}
```

### Check container status:
```bash
docker ps
```

You should see the `altoai-mvp` container running.

---

## Troubleshooting

### Port 8080 already in use:
```bash
# Find what's using port 8080
lsof -i :8080

# Or use a different port in docker-compose.yml:
ports:
  - "3000:8080"  # Change 3000 to any available port
```

### Container won't start:
```bash
# Check logs for errors
docker logs altoai-mvp

# Or with docker-compose
docker-compose logs
```

### Rebuild after code changes:
```bash
# Stop and remove old containers
docker-compose down

# Rebuild and start
docker-compose up --build
```

### Clear everything and start fresh:
```bash
# Stop containers
docker-compose down

# Remove images
docker rmi altoai-mvp:latest

# Rebuild from scratch
docker-compose up --build
```

---

## Environment Variables

Make sure your `.env` file is in the root directory with all required variables:
- Database connection (if using PostgreSQL)
- Google OAuth credentials
- JWT secret
- Other API keys

The `.env` file is automatically loaded by docker-compose.

---

## Development vs Production

**For Development:**
- Use `docker-compose up` to see logs in real-time
- Rebuild after code changes: `docker-compose up --build`

**For Production:**
- Use `docker-compose up -d` to run in background
- Set `restart: unless-stopped` (already configured)
- Monitor with: `docker-compose logs -f`



