# マルチステージビルド: ビルドステージ
FROM node:18-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# package.jsonをコピー
COPY package*.json ./

# 依存関係をインストール
RUN npm ci --only=production

# ソースコードをコピー
COPY . .

# Next.jsアプリケーションをビルド
RUN npm run build

# 本番ステージ: nginx
FROM nginx:alpine

# nginxの設定ファイルを削除
RUN rm /etc/nginx/conf.d/default.conf

# ビルドされたNext.jsアプリをnginxの静的ファイルディレクトリにコピー
COPY --from=builder /app/out /usr/share/nginx/html

# ポート80を公開
EXPOSE 80

# nginxをフォアグラウンドで実行
CMD ["nginx", "-g", "daemon off;"] 