# マルチステージビルド: ビルドステージ
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# Go modulesファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app ./cmd/main.go

# 実行ステージ: 軽量なalpineベース
FROM alpine:latest

# セキュリティアップデートとCA証明書
RUN apk --no-cache add ca-certificates tzdata
RUN update-ca-certificates

# 非rootユーザーを作成
RUN adduser -D -s /bin/sh appuser

# 作業ディレクトリを設定
WORKDIR /app

# ビルドされたバイナリをコピー
COPY --from=builder /app/app .

# バイナリに実行権限を付与
RUN chmod +x ./app

# 非rootユーザーに変更
USER appuser

# ポートを公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションを起動
CMD ["./app"] 