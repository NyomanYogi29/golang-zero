dev:
	air

test:
	go test -v ./internal/...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

build:
	go build -o ./bin/server ./cmd/server.go

push:
	@if [ -z "$(m)" ]; then \
		echo "Error: Pesan commit wajib diisi!"; \
		echo "Contoh: make push m=\"feat: add user registration endpoint\""; \
		exit 1; \
	fi
	git add .
	git commit -m "$(m)"
	git push origin main