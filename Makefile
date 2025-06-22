.PHONY: help setup dev build test lint clean db-migrate db-migrate-down db-migrate-version db-create-migration db-seed db-reset gen-api version

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  setup               - 初期セットアップ"
	@echo "  dev                 - 開発サーバー起動"
	@echo "  prod                - 本番環境起動（nginx）"
	@echo "  build               - プロダクションビルド"
	@echo "  test                - テスト実行"
	@echo "  lint                - リント実行"
	@echo "  db-migrate          - データベースマイグレーション（up）"
	@echo "  db-migrate-down     - マイグレーションロールバック（down）"
	@echo "  db-migrate-version  - マイグレーション状態確認"
	@echo "  db-create-migration - 新しいマイグレーションファイル作成"
	@echo "  db-seed             - テストデータ投入"
	@echo "  db-reset            - データベースリセット"
	@echo "  merge-api           - OpenAPI仕様をマージ"
	@echo "  gen-api             - OpenAPIからコード生成"
	@echo "  version             - バージョン情報表示"
	@echo "  clean               - クリーンアップ"

# 初期セットアップ
setup:
	@echo "🚀 初期セットアップを開始..."
	cp .env.example .env || true
	docker-compose pull
	cd frontend && npm install
	cd backend && go mod download
	@echo "✅ セットアップ完了!"

# 開発環境起動
dev:
	@echo "🔥 開発サーバーを起動中..."
	docker-compose up -d mysql
	sleep 10
	$(MAKE) db-migrate
	docker-compose --profile dev up --build

# 本番環境起動
prod:
	@echo "🚀 本番環境を起動中..."
	docker-compose up -d mysql redis
	sleep 10
	$(MAKE) db-migrate
	docker-compose --profile prod up --build -d
	@echo "✅ 本番環境起動完了! http://localhost でアクセス可能です"

# プロダクションビルド
build:
	@echo "📦 プロダクションビルドを開始..."
	cd frontend && npm run build
	cd backend && go build -o bin/app cmd/main.go
	@echo "✅ ビルド完了!"

# テスト実行（テストファイルがある場合のみ）
test:
	@echo "🧪 テストをチェック中..."
	@if [ -f "frontend/package.json" ] && grep -q '"test"' frontend/package.json; then \
		echo "フロントエンドテストを実行中..."; \
		cd frontend && npm test; \
	else \
		echo "⏭️ フロントエンドテストファイルが見つかりません"; \
	fi
	@if find backend -name "*_test.go" -type f | grep -q .; then \
		echo "バックエンドテストを実行中..."; \
		cd backend && go test ./...; \
	else \
		echo "⏭️ バックエンドテストファイルが見つかりません"; \
	fi
	@echo "✅ テストチェック完了!"

# リント実行
lint:
	@echo "🔍 コード品質をチェック中..."
	cd frontend && npm run lint
	cd backend && golangci-lint run
	@echo "✅ リント完了!"

# データベースマイグレーション
db-migrate:
	@echo "📊 データベースマイグレーションを実行中..."
	docker-compose exec -T mysql mysql -uroot -ppassword -P 3306 -e "CREATE DATABASE IF NOT EXISTS app_db;"
	cd backend && go run cmd/migrate/main.go up
	@echo "✅ マイグレーション完了!"

# マイグレーション（ダウン）
db-migrate-down:
	@echo "📊 マイグレーションをロールバック中..."
	cd backend && go run cmd/migrate/main.go down
	@echo "✅ マイグレーションロールバック完了!"

# マイグレーション状態確認
db-migrate-version:
	@echo "📊 マイグレーション状態を確認中..."
	cd backend && go run cmd/migrate/main.go version
	@echo "✅ マイグレーション状態確認完了!"

# 新しいマイグレーションファイル作成
db-create-migration:
	@echo "📊 新しいマイグレーションファイルを作成中..."
	@read -p "マイグレーション名を入力してください: " name; \
	existing_count=$$(ls database/migrations/*.up.sql 2>/dev/null | wc -l); \
	next_number=$$(printf "%06d" $$((existing_count + 1))); \
	touch database/migrations/$$next_number\_$$name.up.sql; \
	touch database/migrations/$$next_number\_$$name.down.sql; \
	echo "✅ マイグレーションファイル作成完了: $$next_number\_$$name"; \
	echo "UP ファイル: database/migrations/$$next_number\_$$name.up.sql"; \
	echo "DOWN ファイル: database/migrations/$$next_number\_$$name.down.sql"

# テストデータ投入
db-seed:
	@echo "🌱 テストデータを投入中..."
	cd backend && go run cmd/seed/main.go
	@echo "✅ シード完了!"

# データベースリセット
db-reset:
	@echo "🔄 データベースをリセット中..."
	docker-compose down
	docker volume rm $$(docker volume ls -q | grep mysql) || true
	docker-compose up -d mysql
	sleep 10
	$(MAKE) db-migrate
	$(MAKE) db-seed
	@echo "✅ データベースリセット完了!"

# OpenAPI仕様をマージ
merge-api:
	@echo "📋 OpenAPI仕様をマージ中..."
	cd api && yarn install && yarn merge
	@echo "✅ API仕様マージ完了!"

# OpenAPIからコード生成
gen-api: merge-api
	@echo "⚙️ APIコードを生成中..."
	# バックエンド（Go）のコード生成
	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yml \
		-g go-gin-server \
		-o /local/backend/api
	# フロントエンド（TypeScript）のコード生成
	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yml \
		-g typescript-axios \
		-o /local/frontend/src/api
	@echo "✅ コード生成完了!"

# クリーンアップ
clean:
	@echo "🧹 クリーンアップ中..."
	docker-compose --profile dev --profile prod down
	docker system prune -f
	docker volume prune -f
	cd frontend && rm -rf .next out node_modules || true
	cd backend && rm -rf bin || true
	@echo "✅ クリーンアップ完了!"

# バージョン情報表示
version:
	@echo "📋 バージョン情報:"
	@echo "Go: $$(go version)"
	@echo "Node.js: $$(node --version)"
	@echo "Docker: $$(docker --version)"
	@echo "Make: $$(make --version | head -1)"

# 本番デプロイ
deploy:
	@echo "🚀 本番環境へデプロイ中..."
	cd infrastructure && terraform apply -auto-approve
	@echo "✅ デプロイ完了!" 