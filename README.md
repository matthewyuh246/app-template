# Web Application Template

このテンプレートは、モダンな技術スタックを使用してWebアプリケーションを素早く構築するための再利用可能なテンプレートです。

## 技術スタック

### Frontend
- **Next.js 14** (TypeScript)
- **TailwindCSS** - スタイリング
- **React Hook Form** - フォーム管理
- **React Query** - データフェッチング
- **Nginx** - リバースプロキシ・静的ファイル配信

### Backend
- **Gin** (Go) - Webフレームワーク
- **SQLBoiler** - ORM
- **Swagger/OpenAPI** - API仕様とコード生成
- **JWT** - 認証

### Database
- **MySQL 8.0**

### Infrastructure
- **Docker & Docker Compose** - 開発環境
- **GitHub Actions** - CI/CD
- **AWS** - クラウドインフラ
- **Terraform** - インフラストラクチャ as Code

### Architecture
- **Clean Architecture**
  - Repository層 - データアクセス
  - UseCase層 - ビジネスロジック
  - Controller層 - リクエストハンドリング

### 開発ツール
- **Make** - ビルドオートメーション
- **Air** - Go のホットリロード
- **ESLint/Prettier** - コード品質
- **golangci-lint** - Go静的解析

## CI/CD パイプライン

### ブランチ戦略

```
📦 Repository
├── 🔒 main (本番環境) - develop/からのみマージ可能
│   └── 🚀 自動デプロイ → AWS ECS
│
├── 🔧 develop/** (開発ブランチ)
│   ├── 🚀 自動CI (Lint + Build)
│   └── 👥 コードレビュー必須
│
├── feature/** 
├── fix/**
└── hotfix/**
```

### 自動CI（高速フィードバック）

**実行タイミング:**
- Pull Request作成・更新時
- develop/**, feature/**, fix/**, hotfix/** ブランチへのpush時

**チェック内容:**
1. **ブランチ名検証** - mainへのPRは`develop/`ブランチからのみ
2. **変更ファイル検出** - backend/frontend の変更を検出
3. **Backend CI** (変更時のみ実行)
   - Go 1.24 環境でビルド
   - golangci-lint による静的解析
4. **Frontend CI** (変更時のみ実行)
   - Node.js 18 環境でビルド
   - ESLint による静的解析
   - TypeScript型チェック

**実行時間**: 約1-3分

### デプロイパイプライン

**実行タイミング:**
- mainブランチへのpush時
- 手動実行 (workflow_dispatch)

**パイプライン:**
1. **Pre-deploy checks** - develop/ブランチからのマージ確認
2. **Build** - 本番用ビルド
3. **Deploy** - AWS ECS/ECRへのデプロイ
4. **Smoke Tests** - 本番環境での動作確認
5. **Rollback** - 失敗時の自動ロールバック

## プロジェクト構造

```
.
├── frontend/                 # Next.js アプリケーション
├── backend/                  # Go API サーバー
├── nginx/                    # Nginx設定
├── database/                 # データベース関連
├── infrastructure/           # Terraform設定
├── api/                      # OpenAPI仕様
├── .github/workflows/        # CI/CD設定
├── docker-compose.yml        # 開発・本番環境
├── Makefile                  # ビルドタスク
└── README.md
```

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.24+
- Node.js 18+
- Make

### 初期セットアップ

1. **リポジトリのクローン**
```bash
git clone <your-repo>
cd App-Template
```

2. **環境変数の設定**
```bash
cp .env.example .env
# .envファイルを編集して適切な値を設定
```

3. **依存関係のインストールと初期化**
```bash
make setup
```

4. **開発サーバーの起動**
```bash
# 開発環境（Next.js開発サーバー）
make dev

# 本番環境（Nginx + 静的ビルド）
make prod
```

### 利用可能なMakeコマンド

```bash
make setup          # 初期セットアップ
make dev             # 開発サーバー起動
make build           # プロダクションビルド
make test            # テスト実行
make lint            # リント実行
make db-migrate      # データベースマイグレーション
make db-seed         # テストデータ投入
make gen-api         # OpenAPIからコード生成
make clean           # クリーンアップ
```

## 開発フロー

### 1. 新機能開発

```bash
# develop/ブランチ作成
git checkout -b develop/feature-name

# 開発・テスト
make dev
make lint
make test

# Pull Request作成
git push origin develop/feature-name
# GitHub UI でPR作成
```

### 2. 自動CI実行

- PR作成時に自動でCI実行
- 変更されたコンポーネント（backend/frontend）のみテスト
- 1-3分で結果がフィードバック

### 3. コードレビューとマージ

- コードレビュー実施
- CI通過後にmainブランチへマージ
- 自動デプロイが実行

### 4. デプロイ

- mainブランチマージ時に自動実行
- AWS ECS/ECRへのデプロイ
- スモークテスト実行
- 失敗時は自動ロールバック

## API仕様

### 構造化されたAPI定義

可読性向上のため、API仕様は以下のように分割されています：

```
api/
├── package.json              # yarn merge設定
├── index.yml                 # メイン設定ファイル
├── openapi.yml              # 統合済みファイル（自動生成）
├── info/
│   └── info.yml            # API基本情報
├── paths/
│   ├── health.yml          # ヘルスチェックエンドポイント
│   ├── auth.yml           # 認証エンドポイント
│   └── users.yml          # ユーザーエンドポイント
└── components/
    ├── schemas/
    │   ├── user.yml       # ユーザー関連スキーマ
    │   ├── auth.yml       # 認証関連スキーマ
    │   ├── common.yml     # 共通スキーマ
    │   └── error.yml      # エラースキーマ
    └── security/
        └── security.yml   # セキュリティスキーム
```

### API仕様の使用方法

```bash
# 分割されたファイルを統合
make merge-api

# API仕様を統合してコード生成
make gen-api
```

統合されたAPI仕様は以下でアクセス可能です：
- Swagger UI: http://localhost:8080/swagger/
- API仕様: http://localhost:8080/swagger/doc.json

## データベース

### マイグレーション
```bash
make db-migrate
```

### テストデータ
```bash
make db-seed
```

## デプロイ

### AWS環境へのデプロイ
```bash
cd infrastructure
terraform init
terraform plan
terraform apply
```

### GitHub Actions
- プルリクエスト時：テスト、リント実行
- mainブランチマージ時：ビルド、デプロイ

## 開発フロー

1. **新機能開発**
   - `api/` ディレクトリ内の分割されたファイルでAPI仕様を定義
   - `make merge-api` で仕様を統合
   - `make gen-api` でコード生成
   - バックエンド実装 (Repository → UseCase → Controller)
   - フロントエンド実装

2. **テスト**
   - `make test` でテスト実行
   - `make lint` でコード品質チェック

3. **デプロイ**
   - プルリクエスト作成
   - CI/CDパイプライン実行
   - レビュー後マージ

## アーキテクチャ

### 環境構成

**開発環境 (`make dev`)**
- Next.js開発サーバー (http://localhost:3000)
- Go APIサーバー (http://localhost:8080)
- MySQL, Redis
- Swagger UI (http://localhost:8081)

**本番環境 (`make prod`)**
- Nginx リバースプロキシ (http://localhost:80)
  - 静的ファイル配信 (Next.js静的エクスポート)
  - `/api/*` をGo APIサーバーにプロキシ
  - gzip圧縮、キャッシング、セキュリティヘッダー
- Go APIサーバー
- MySQL, Redis

### Nginx構成

```nginx
# 静的ファイル (Next.js)
location / {
    try_files $uri $uri/ @nextjs;
}

# APIプロキシ
location /api/ {
    proxy_pass http://backend/;
}

# ヘルスチェック
location /health {
    proxy_pass http://backend/health;
}
```

### Clean Architecture層構造

```
┌─────────────────┐
│   Controller    │ ← HTTP リクエスト処理
├─────────────────┤
│    UseCase      │ ← ビジネスロジック
├─────────────────┤
│   Repository    │ ← データアクセス
└─────────────────┘
```

### ディレクトリ構造 (Backend)
```
backend/
├── cmd/                     # エントリーポイント
├── internal/
│   ├── controller/          # HTTPハンドラー
│   ├── usecase/            # ビジネスロジック
│   ├── repository/         # データアクセス
│   ├── entity/             # エンティティ
│   └── infrastructure/     # 外部依存
├── api/                    # 生成されたAPIコード
└── pkg/                    # 共通パッケージ
```

## カスタマイズ

このテンプレートをプロジェクトに合わせてカスタマイズ：

1. **プロジェクト名の変更**
   - `go.mod` のモジュール名
   - `package.json` のプロジェクト名
   - Docker設定

2. **API仕様の修正**
   - `api/` ディレクトリ内の適切なファイルを編集
   - `make merge-api` で統合
   - `make gen-api` で再生成

3. **データベーススキーマ**
   - `database/migrations/` でマイグレーション追加
   - SQLBoilerモデル再生成

## トラブルシューティング

### よくある問題

1. **ポートが使用中**
```bash
make clean
docker-compose down
```

2. **データベース接続エラー**
```bash
make db-reset
```

3. **依存関係の問題**
```bash
make setup
```

## ライセンス

MIT License

## 貢献

1. `develop/feature-name` ブランチを作成
2. 変更を実装
3. CI通過を確認
4. Pull Request作成
5. コードレビュー
6. mainブランチへマージ

---

質問や問題がある場合は、Issueを作成してください。 