FROM golang:1.21-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /calendar-api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Tokyo

WORKDIR /root/

COPY --from=builder /calendar-api .

EXPOSE 8080

CMD ["./calendar-api"]
