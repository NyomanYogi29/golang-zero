dev:
	air

test:
	go test -v ./internal/...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

swagger:
	swag init -g cmd/server.go -o docs

build: swagger
	go build -o ./bin/server ./cmd/server.go

tools:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest

push:
	@if [ -z "$(m)" ]; then \
		echo "Error: Pesan commit wajib diisi!"; \
		echo "Contoh: make push m=\"feat: add user registration endpoint\""; \
		exit 1; \
	fi
	git add .
	git commit -m "$(m)"
	git push origin main