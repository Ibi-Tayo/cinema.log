# Stop Angular app (running on default port 4200)
lsof -ti:4200 | xargs kill -9 2>/dev/null

# Stop Go API server (adjust port if different, assuming 8080)
lsof -ti:8080 | xargs kill -9 2>/dev/null

# Stop Docker containers
cd cinema.log.server.golang/
docker compose -f 'docker-compose.yml' down

# Optional: Stop Docker Desktop
osascript -e 'tell application "Docker" to quit'

echo "Development environment stopped"