FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o chat-service .

# Build the final image
FROM amd64/alpine:latest

# Install necessary dependencies
RUN apk add --no-cache libc6-compat

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/chat-service .

# Verify if the binary exists in the final image
RUN ls -l /root/

CMD ["./chat-service"]

