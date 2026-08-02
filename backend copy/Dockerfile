FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /out/api \
    ./cmd/api

FROM alpine:3.20 AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Jakarta

RUN mkdir -p /app/public/uploads

COPY --from=builder /out/api /app/api

EXPOSE 8080

CMD ["/app/api"]
