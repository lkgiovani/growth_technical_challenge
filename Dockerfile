ARG GO_VERSION=1.24

# Build
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /service
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./cmd/server

# Test
FROM golang:${GO_VERSION}-alpine AS tests
ENV CI=1
WORKDIR /service
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go clean -testcache
RUN go test -v ./...

# Docs Build
FROM node:22-slim AS docs
WORKDIR /docs
COPY ./docs/package.json ./docs/package-lock.json ./
RUN npm install
COPY ./docs/ ./
RUN npm run compile

# Image
FROM gcr.io/distroless/static-debian12 AS production
ENV PROFILE=prod
WORKDIR /service
USER nonroot:nonroot
COPY --from=docs /docs/schema ./docs/schema
COPY --from=build --chown=nonroot:nonroot /app ./app
ENTRYPOINT ["/service/app"]