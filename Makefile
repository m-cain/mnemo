.PHONY: all db backend frontend clean

all: db backend frontend

db:
	@echo "Starting database..."
	docker-compose up -d db

backend:
	@echo "Starting backend with Air..."
	cd backend && air

frontend:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

clean:
	@echo "Stopping and removing containers..."
	docker-compose down
