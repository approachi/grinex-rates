.PHONY: build
build:
	go build -o bin/server ./cmd/server

.PHONY: test-integration
test-integration:
	docker-compose -f docker-compose.test.yml up -d
	docker-compose -f docker-compose.test.yml wait migrator-test
	go test -count=1 -tags integration ./tests/integration/... -v; \
	docker-compose -f docker-compose.test.yml down

.PHONY: test
test:
	go test -count=1 ./internal/... -v

.PHONY: ftest
ftest: test test-integration

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
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run ./...

.PHONY: proto
proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/approachi/grinex-rates \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/approachi/grinex-rates \
		proto/rates/rates.proto
