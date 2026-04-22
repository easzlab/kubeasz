.PHONY: dev-up dev-down migrate-up migrate-down build-backend build-frontend

# Development targets
dev-up:
	docker-compose up -d --build

# Wait for services to be healthy, then show status
dev-ready:
	@echo "Waiting for services..."
	@sleep 5
	@docker-compose ps
	@echo ""
	@echo "MySQL: mysql -h 127.0.0.1 -P 3306 -u ksk8s -p"
	@echo "Backend API: http://localhost:8080"
	@echo "Frontend: http://localhost"

dev-down:
	docker-compose down

dev-down-volumes:
	docker-compose down -v

migrate-up:
	cd ksk8s-backend && go run cmd/migrate/main.go up

migrate-down:
	cd ksk8s-backend && go run cmd/migrate/main.go down

build-backend:
	docker build -t ksk8s-backend:latest ./ksk8s-backend

build-frontend:
	docker build -t ksk8s-frontend:latest ./ksk8s-frontend
