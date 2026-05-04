# syntax=docker/dockerfile:1

FROM golang:1.24 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./backend

FROM gcr.io/distroless/static-debian12
COPY --from=build /backend /backend
EXPOSE 8080
CMD ["/backend"]
