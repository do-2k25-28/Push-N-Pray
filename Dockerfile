FROM golang:1.26.2 AS builder

WORKDIR /src

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian13

USER nonroot:nonroot

COPY --from=build /app/push-n-pray /usr/local/bin/push-n-pray

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/push-n-pray"]
