.PHONY: test vet run docker-up

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/gmf-core

docker-up:
	docker compose up --build
