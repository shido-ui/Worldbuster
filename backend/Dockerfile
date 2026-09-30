FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY backend ./backend
RUN CGO_ENABLED=0 go build -o /out/worldbuster ./backend/cmd/worldbuster

FROM alpine:3.22
WORKDIR /app
COPY --from=build /out/worldbuster /app/worldbuster
COPY database/migrations /app/database/migrations
EXPOSE 8080
ENV WORLDBUSTER_MIGRATIONS_DIR=/app/database/migrations
ENTRYPOINT ["/app/worldbuster"]
