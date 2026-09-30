FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY backend ./backend
RUN CGO_ENABLED=0 go build -o /out/worldbuster ./backend/cmd/worldbuster
FROM alpine:3.22
COPY --from=build /out/worldbuster /worldbuster
EXPOSE 8080
ENTRYPOINT ["/worldbuster"]
