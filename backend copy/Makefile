# ==================================================
# Schoolpay - BACKEND MAKEFILE
# ==================================================
# Dipakai dari dalam backend/ atau via root Makefile.
# ==================================================

SHELL := /bin/sh

ENV_FILE ?= ../.env
MIGRATIONS_DIR := ./migrations

# Load root .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

# Goose command (deferred expansion biar .env terbaca dulu)
GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) mysql '$(DB_USER):$(DB_PASS)@tcp(localhost:$(DB_PORT_EXTERNAL))/$(DB_NAME)?parseTime=true&multiStatements=true&charset=utf8mb4&loc=Local'

# ==================================================
# HELP
# ==================================================
.PHONY: help
help:
	@echo ""
	@echo "Schoolpay Backend"
	@echo "=================================================="
	@echo "Build & Run:"
	@echo "  make build                - Build semua package Go"
	@echo "  make run                  - Jalankan API server"
	@echo "  make worker               - Jalankan background worker"
	@echo "  make seed                 - Isi data awal"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt                  - Format kode Go"
	@echo "  make lint                 - GolangCI-Lint"
	@echo "  make test                 - Test backend"
	@echo "  make test-race            - Test dengan race detection"
	@echo ""
	@echo "Migration:"
	@echo "  make migrate-create name=...  - Buat file migrasi baru"
	@echo "  make migrate-status           - Lihat status migrasi"
	@echo "  make migrate-up               - Jalankan migrasi pending"
	@echo "  make migrate-down             - Rollback 1 migrasi"
	@echo "  make migrate-redo             - Rollback + jalankan ulang"
	@echo "  make migrate-reset            - Rollback semua"
	@echo "  make migrate-version          - Versi migrasi saat ini"
	@echo ""

# ==================================================
# BUILD & RUN
# ==================================================
.PHONY: build
build:
	go build -v ./...

.PHONY: run
run:
	go run ./cmd/api

.PHONY: run-local
run-local:
	DB_HOST=127.0.0.1 DB_PORT=$(DB_PORT_EXTERNAL) REDIS_HOST=127.0.0.1 REDIS_PORT=$(REDIS_PORT_EXTERNAL) go run ./cmd/api

.PHONY: worker
worker:
	go run ./cmd/worker

.PHONY: seed
seed:
	go run ./cmd/seeder

# ==================================================
# CODE QUALITY
# ==================================================
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

# ==================================================
# MIGRATION
# ==================================================
.PHONY: migrate-create
migrate-create:
	@test -n "$(name)" || (echo "Isi nama migration. Contoh: make migrate-create name=create_users" && exit 1)
	@mkdir -p "$(MIGRATIONS_DIR)"
	goose -dir $(MIGRATIONS_DIR) create "$(name)" sql

.PHONY: migrate-status
migrate-status:
	$(GOOSE_CMD) status

.PHONY: migrate-up
migrate-up:
	$(GOOSE_CMD) up

.PHONY: migrate-down
migrate-down:
	$(GOOSE_CMD) down

.PHONY: migrate-redo
migrate-redo:
	$(GOOSE_CMD) redo

.PHONY: migrate-reset
migrate-reset:
	$(GOOSE_CMD) reset

.PHONY: migrate-version
migrate-version:
	$(GOOSE_CMD) version
