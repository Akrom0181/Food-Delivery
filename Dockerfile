FROM golang:1.23 AS builder

WORKDIR /app

COPY serviceAccountKey.json ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Debug: Check contents of ./cmd/app
RUN echo "Checking ./cmd/app contents:" && ls -l /app/cmd/app

RUN GOOS=linux GOARCH=amd64 go build -tags migrate -o food ./cmd/app

# Debug: Verify binary exists
RUN echo "Checking binary:" && ls -l /app/food

# Step 2: Create a lightweight image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/serviceAccountKey.json ./
COPY --from=builder /app/food /app/
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

RUN apk add --no-cache libc6-compat

EXPOSE 9090

CMD ["/app/food"]
