up:
	cd deployments && docker compose up --build -d
down:
	cd deployments && docker compose down
ps:
	cd deployments && docker compose ps
logs:
	cd deployments && docker compose logs -f

test-unit:
	go test ./... -short
test-integration:
	go test -tags=integration ./test/integration/... -v

lint:
	golangci-lint run ./...