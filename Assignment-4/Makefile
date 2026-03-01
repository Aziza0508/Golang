.PHONY: up down restart logs ps images build scale clean

# Build and start all services
up:
	docker compose up --build -d

# Stop and remove containers (preserves volumes)
down:
	docker compose down

# Restart all services
restart:
	docker compose down
	docker compose up --build -d

# View logs (follow mode)
logs:
	docker compose logs -f

# Show running containers
ps:
	docker compose ps

# Show docker images and sizes
images:
	docker images

# Build without starting
build:
	docker compose build

# Scale web-app to 3 instances with nginx load balancing
scale:
	docker compose up --build --scale web-app=3 -d

# Stop, remove containers, volumes, and images
clean:
	docker compose down -v --rmi local
