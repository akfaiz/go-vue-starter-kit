SHELL := /bin/bash
APP := yourapp
BIN := bin/$(APP)

# Default env (override with `make run DEV=1`)
DEV ?= 0

.PHONY: all tidy run dev build build-embed web-install web-build clean

all: build-embed

## ---- Go targets ----

tidy:
	@echo ">> go mod tidy"
	go mod tidy

run:
	@echo ">> Running server (DEV=$(DEV))"
	DEV=$(DEV) go run . serve

run-dev:
	@echo ">> Running server with air (auto-reload on code changes)"
	DEV=1 air -c .air.toml

build:
	@echo ">> Building Go binary only"
	go build -o $(BIN) .

build-embed: web-build
	@echo ">> Building Go binary with embedded Vue dist"
	go build -tags=netgo -ldflags="-s -w" -o $(BIN) .

## ---- Vue / UI targets ----

web-install:
	@echo ">> Installing UI deps with pnpm"
	cd ui && pnpm install

web-build: web-install
	@echo ">> Building Vue app into /web/dist"
	cd ui && pnpm build

run-web-dev:
	@echo ">> Starting Vite dev server (port 5173)"
	cd ui && pnpm dev

## ---- Dev workflow ----

dev:
	@echo ">> Start frontend dev server (Vite)"
	@echo "   Run in one terminal: cd ui && pnpm dev"
	@echo "   Run in another: DEV=1 go run . serve"
	@echo "   Browse http://localhost:8080 (Echo → proxy → Vite)"

## ---- Cleanup ----

clean:
	@echo ">> Cleaning build artifacts"
	rm -rf $(BIN) web/dist
	cd ui && rm -rf node_modules .vite

docker-build:
	docker build --build-arg GIT_COMMIT="$(git rev-parse --short HEAD)" --build-arg BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)" -t go-vue-app .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env go-vue-app serve 