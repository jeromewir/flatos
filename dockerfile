# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . .
ENV GIN_MODE=release
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go
CMD ["/app/start.sh"]

# # Final stage
# FROM scratch
# WORKDIR /app
# ENV GIN_MODE=release
# COPY --from=builder /app/main .
# COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
# COPY --from=builder /app/migrations ./migrations
# COPY --from=builder /app/makefile ./makefile
# EXPOSE 8080
# ENTRYPOINT ["/app/main"]
