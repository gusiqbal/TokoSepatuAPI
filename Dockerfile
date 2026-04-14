FROM golang:1.25-alpine AS builder

# Enable CGO only if you need C bindings
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /root
COPY --from=builder /app/main .
# COPY .env .env bahaya env bisa bocor ke public lewat image

EXPOSE 8080
CMD ["./main"]
