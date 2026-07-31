# ==================================================
# DPP Gradasi - ROOT MAKEFILE (Local Migration)
# ==================================================

SHELL := /bin/sh

# Default environment file
ENV_FILE ?= .env

# Load .env file
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

# Default database variables if not set in .env
DB_USER ?= root
DB_PASS ?= 
DB_HOST ?= 127.0.0.1
DB_PORT ?= 3306
DB_NAME ?= dpp_gradasi

MIGRATIONS_DIR := ./backend/migrations

# Goose command configuration for local MySQL
GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) mysql '$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&multiStatements=true&charset=utf8mb4&loc=Local'

.PHONY: help
help:
	@echo "=================================================="
	@echo "DPP Gradasi Commands"
	@echo "=================================================="
	@echo "Build & Run (Backend):"
	@echo "  make build                - Build semua package Go"
	@echo "  make run                  - Jalankan API server"
	@echo "  make run-local            - Jalankan API server (lokal DB config)"
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
	@echo "  make migrate-reset            - Rollback semua migrasi"
	@echo "  make migrate-version          - Versi migrasi saat ini"
	@echo "=================================================="

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

# ==================================================
# BUILD & RUN (Backend)
# ==================================================
.PHONY: build
build:
	cd backend && go build -v ./...

.PHONY: run
run:
	cd backend && go run ./cmd/api

.PHONY: run-local
run-local:
	cd backend && DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASS=$(DB_PASS) DB_NAME=$(DB_NAME) go run ./cmd/api

.PHONY: worker
worker:
	cd backend && go run ./cmd/worker

.PHONY: seed
seed:
	cd backend && go run ./cmd/seeder

# ==================================================
# CODE QUALITY
# ==================================================
.PHONY: fmt
fmt:
	cd backend && go fmt ./...

.PHONY: lint
lint:
	cd backend && golangci-lint run

.PHONY: test
test:
	cd backend && go test ./...

.PHONY: test-race
test-race:
	cd backend && go test -race ./...
