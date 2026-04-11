# syntax=docker/dockerfile:1.7

# ---------- build ----------
FROM golang:1.23-alpine AS build
WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum* ./

# If go.sum isn't committed yet, tidy will fetch deps + generate it.
COPY . .
RUN go mod tidy

# Static build, stripped, trimmed for small image
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/api ./cmd/api

# ---------- runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /app

COPY --from=build /out/api /app/api
COPY migrations /app/migrations
COPY web /app/web

ENV HTTP_ADDR=":8182" \
    MIGRATIONS_DIR="/app/migrations" \
    TEMPLATES_DIR="/app/web/templates" \
    APP_ENV="production"

EXPOSE 8182
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
