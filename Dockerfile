# --- build stage ---
FROM golang:1.26-alpine AS builder

ARG SERVICE

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${SERVICE}

# --- run stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/app /app/app

CMD ["/app/app"]