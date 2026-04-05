.PHONY: build
build:
	go build -o bin/server ./cmd/server

.PHONY: test-unit
test-unit:
	go test -count=1 ./internal/... -v

.PHONY: test-integration
test-integration:
	docker-compose -f docker-compose.test.yml up -d
	docker-compose -f docker-compose.test.yml wait migrator-test
	go test -count=1 ./tests/integration/... -v; \
	docker-compose -f docker-compose.test.yml down

.PHONY: test
test: test-unit test-integration

.PHONY: docker-build
docker-build:
	docker build -t grinex-rates .

.PHONY: run
run:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: db-conn
db-conn:
	docker-compose exec postgres psql -U postgres -d grinex

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: proto
proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/approachi/grinex-rates \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/approachi/grinex-rates \
		proto/rates/rates.proto
