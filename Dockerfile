# ---------------------------------------
# 1) Frontend build (Vue/Vite with PNPM)
# ---------------------------------------
FROM node:20-alpine AS frontend
RUN corepack enable && apk add --no-cache git
WORKDIR /app/ui

# Install deps with caching
COPY ui/ .
RUN --mount=type=cache,target=/root/.pnpm-store pnpm install --frozen-lockfile

# Build frontend
RUN --mount=type=cache,target=/root/.pnpm-store pnpm build

# ---------------------------------------
# 2) Backend build (Go at repo root)
# ---------------------------------------
FROM golang:1.25-alpine AS backend
RUN apk add --no-cache git build-base
WORKDIR /app

# Go deps first for cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy backend source
COPY . .

# Copy built frontend dist into /web/dist
RUN rm -rf ./web/dist && mkdir -p ./web/dist
COPY --from=frontend /app/web/dist ./web/dist

# Build Go binary (embed will include ./web/dist)
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w -X main.Commit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" -o /app/bin/app .

# ---------------------------------------
# 3) Runtime (minimal)
# ---------------------------------------
FROM alpine:3.20 AS runtime
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=backend /app/bin/app /app/app

EXPOSE 8080

USER app
ENTRYPOINT ["/app/app"]