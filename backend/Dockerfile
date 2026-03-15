FROM golang:1.26-trixie AS builder

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /calendar-api ./cmd/api

FROM debian:trixie-slim

RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
ENV TZ=Asia/Tokyo

WORKDIR /root/

COPY --from=builder /calendar-api .

EXPOSE 8080

CMD ["./calendar-api"]
