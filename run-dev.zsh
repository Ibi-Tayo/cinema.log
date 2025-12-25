# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLIENT_DIR="$SCRIPT_DIR/cinema.log.client"

# Start Angular App
osascript -e 'tell application "Terminal" to do script "cd '"$CLIENT_DIR"' && npm run start"'

# Make sure docker is running
open -a Docker

# Wait for Docker to be ready
echo "Waiting for Docker to start..."
while ! docker info >/dev/null 2>&1; do
    sleep 1
done
echo "Docker is ready!"

# Spin up container
cd cinema.log.server.golang/
docker compose -f 'docker-compose.yml' down
docker compose -f 'docker-compose.yml' up -d 'psql_bp' 

# Wait for postgres to be healthy
echo "Waiting for PostgreSQL to be ready..."
sleep 2

# Start server
go run cmd/api/main.go 

