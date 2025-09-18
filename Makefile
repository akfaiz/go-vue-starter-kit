SHELL := /bin/bash
APP   := govue
BIN   := bin/$(APP)

DEV        ?= 0
PKG        ?= ./...
PROCS      ?= 4
TIMEOUT    ?= 5m
RACE       ?=
TAGS       ?=
COVERFILE  ?= coverage.out
JUNITFILE  ?= junit.xml

all: build-embed

# coverpkg across selected packages
COVERPKG   := $(shell go list $(PKG) | tr '\n' ',')

# ginkgo detection + helper
GINKGO := $(shell command -v ginkgo 2>/dev/null)
ifndef GINKGO
_need_ginkgo = @echo "ginkgo CLI not found. Install: go install github.com/onsi/ginkgo/v2/ginkgo@latest" && exit 2
endif

.PHONY: help
help: ## Show this help
	@printf "\nTargets:\n"
	@grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | \
	awk 'BEGIN { FS=":.*##" } \
	{ gsub(/^[ \t]+|[ \t]+$$/, "", $$1); gsub(/^[ \t]+|[ \t]+$$/, "", $$2); \
	  printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }'

# ---------------------------------------------------------------------

.PHONY: test
test: ## Run all tests (Ginkgo + plain go tests) with coverage
	$(if $(GINKGO),,@$(call _need_ginkgo))
	$(GINKGO) -r -v -procs=$(PROCS) --timeout=$(TIMEOUT) $(RACE) $(TAGS)

.PHONY: test.cover
test.cover: ## Run all tests with coverage (Ginkgo + plain go tests)
	$(if $(GINKGO),,@$(call _need_ginkgo))
	$(GINKGO) -procs=$(PROCS) -cover --randomize-suites \
		--coverprofile=$(COVERFILE) --coverpkg=$(COVERPKG) \
		--trace --timeout=$(TIMEOUT) $(RACE) $(TAGS) $(PKG)
	@echo "Coverage report -> $(COVERFILE)"
	@go tool cover -func=$(COVERFILE) | grep total | awk '{print "Total coverage: "$$3}'

.PHONY: test.ci
test.ci: ## Run tests for CI (Ginkgo + plain go tests) with coverage and JUnit report
	$(if $(GINKGO),,@$(call _need_ginkgo))
	$(GINKGO) -r -v -procs=$(PROCS) --timeout=$(TIMEOUT) $(RACE) $(TAGS) \
		--junit-report=$(JUNITFILE) --coverprofile=$(COVERFILE) --coverpkg=$(COVERPKG) $(PKG)
	@echo "JUnit report -> $(JUNITFILE)"
	@echo "Coverage report -> $(COVERFILE)"
	@go tool cover -func=$(COVERFILE) | grep total | awk '{print "Total coverage: "$$3}'


.PHONY: cover
cover: ## Open HTML coverage from last run
	@go tool cover -html=$(COVERFILE) -o coverage.html
	@echo "Coverage HTML -> coverage.html"
	open coverage.html

.PHONY: tools.ginkgo
tools.ginkgo: ## Install ginkgo CLI
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo ">> go mod tidy"
	go mod tidy

.PHONY: lint
lint: ## Run golangci-lint
	@echo ">> Running golangci-lint"
	golangci-lint run ./...

.PHONY: run
run: ## Run server (DEV=1 for dev mode, with .env)
	@echo ">> Running server (DEV=$(DEV))"
	DEV=$(DEV) go run . serve

.PHONY: run-dev
run-dev: ## Run server with air (auto-reload on code changes)
	@echo ">> Running server with air (auto-reload on code changes)"
	DEV=1 air -c .air.toml

.PHONY: build
build: ## Build Go binary only (no embedded Vue)
	@echo ">> Building Go binary only"
	go build -o $(BIN) .

.PHONY: build-embed
build-embed: web-build ## Build Go binary with embedded Vue dist (requires `make web-build` first)
	@echo ">> Building Go binary with embedded Vue dist"
	go build -tags=netgo -ldflags="-s -w" -o $(BIN) .

.PHONY: web-install
web-install: ## Install UI deps with pnpm (in ui/ folder)
	@echo ">> Installing UI deps with pnpm"
	cd ui && pnpm install

.PHONY: web-build
web-build: web-install ## Build Vue app into web/dist (requires `make web-install` first)
	@echo ">> Building Vue app into /web/dist"
	cd ui && pnpm build

.PHONY: run-web-dev
run-web-dev: ## Run Vite dev server (in ui/ folder)
	@echo ">> Starting Vite dev server (port 5173)"
	cd ui && pnpm dev

.PHONY: clean
clean: ## Clean build artifacts
	@echo ">> Cleaning build artifacts"
	rm -rf $(BIN) web/dist
	cd ui && rm -rf node_modules .vite

.PHONY: docker-build
docker-build: ## Build Docker image with git commit and build time info
	docker build --build-arg GIT_COMMIT="$(git rev-parse --short HEAD)" --build-arg BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)" -t go-vue-app .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 --env-file .env go-vue-app serve