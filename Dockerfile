FROM oven/bun:1-alpine AS frontend-builder
WORKDIR /app/web

COPY web/package.json web/bun.lock* ./
RUN bun install --frozen-lockfile

COPY web ./
RUN bun run build

FROM golang:1.24-alpine AS go-builder
WORKDIR /app

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache gcc musl-dev build-base git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY openapi ./openapi
COPY prompts ./prompts
COPY migrations ./migrations
COPY --from=frontend-builder /app/pkg/web/dist ./pkg/web/dist

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -tags="sqlite_omit_load_extension,osusergo,netgo" \
    -o /app/bin/fincher \
    ./cmd/fincher

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -tags="sqlite_omit_load_extension,osusergo,netgo" \
    -o /app/bin/fincher-seed \
    ./cmd/seed

FROM alpine:3.21 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

RUN addgroup -S -g 10001 fincher && adduser -S -u 10001 -G fincher fincher
COPY --from=go-builder --chown=fincher:fincher /app/bin/fincher /app/fincher
COPY --from=go-builder --chown=fincher:fincher /app/bin/fincher-seed /app/fincher-seed

USER fincher:fincher

ENV FINCHER_ENV=production \
    FINCHER_PORT=8080 \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/fincher"]
