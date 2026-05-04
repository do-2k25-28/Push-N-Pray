# syntax=docker/dockerfile:1

FROM golang:1.24 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY backend ./backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./backend

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /backend /backend
CMD ["/backend"]
