# Swagger/OpenAPI 完全ガイド

## 📋 目次

1. [概要](#概要)
2. [プロジェクト構造](#プロジェクト構造)
3. [基本的な使い方](#基本的な使い方)
4. [API仕様の編集方法](#api仕様の編集方法)
5. [コード生成](#コード生成)
6. [開発フロー](#開発フロー)
7. [実践例](#実践例)
8. [トラブルシューティング](#トラブルシューティング)
9. [ベストプラクティス](#ベストプラクティス)

---

## 概要

このプロジェクトでは、**構造化されたOpenAPI仕様**と**自動コード生成**を組み合わせて、効率的なAPI開発を実現しています。

### 🎯 メリット

- **設計ファースト**: API仕様を先に定義してから実装
- **自動コード生成**: バックエンド・フロントエンドのコードを自動生成
- **ドキュメント自動更新**: Swagger UIで常に最新のAPI仕様を確認可能
- **型安全**: TypeScriptとGoで型安全なAPI通信
- **チーム開発**: 仕様を分割して並行開発可能

### 🏗️ アーキテクチャ

```
分割されたAPI仕様 → swagger-merger → 統合されたopenapi.yml
                                          ↓
                  OpenAPI Generator → Go Server Code + TypeScript Client Code
                                          ↓
                      Swagger UI ← API Documentation
```

---

## プロジェクト構造

### 📁 ディレクトリ構成

```
api/
├── package.json              # npm設定・マージスクリプト
├── index.yml                 # メインAPI仕様ファイル
├── openapi.yml              # 統合済み仕様（自動生成）
├── info/
│   └── info.yml             # API基本情報
├── paths/                   # エンドポイント定義
│   ├── health.yml           # ヘルスチェック
│   ├── auth.yml            # 認証関連
│   └── users.yml           # ユーザー関連
└── components/              # 再利用可能コンポーネント
    ├── schemas/             # データモデル定義
    │   ├── user.yml        # ユーザースキーマ
    │   ├── auth.yml        # 認証スキーマ
    │   ├── common.yml      # 共通スキーマ
    │   └── error.yml       # エラースキーマ
    └── security/
        └── security.yml    # セキュリティスキーム
```

### 🔧 各ファイルの役割

| ファイル | 役割 | 編集頻度 |
|---------|------|----------|
| `index.yml` | 全体設定・参照定義 | 低 |
| `paths/*.yml` | エンドポイント仕様 | 高 |
| `components/schemas/*.yml` | データモデル | 中 |
| `info/info.yml` | API基本情報 | 低 |
| `openapi.yml` | 統合後仕様（自動生成） | 編集禁止 |

---

## 基本的な使い方

### 🚀 クイックスタート

#### 1. 初期セットアップ
```bash
# プロジェクトルートで実行
make setup

# API仕様の初回マージ
make merge-api
```

#### 2. Swagger UI でAPI確認
```bash
# 開発サーバー起動
make dev

# ブラウザで確認
# http://localhost:8080/swagger/
```

#### 3. API仕様編集 → コード生成
```bash
# 1. API仕様を編集（後述）
vim api/paths/users.yml

# 2. 仕様をマージ
make merge-api

# 3. コード生成
make gen-api

# 4. 生成されたコードを確認
ls backend/api/        # Go server code
ls frontend/src/api/   # TypeScript client code
```

### 📖 利用可能なコマンド

```bash
make merge-api    # API仕様をマージ
make gen-api      # コード生成（merge-api も実行）
make dev          # 開発サーバー起動（Swagger UI含む）
```

---

## API仕様の編集方法

### 📝 新しいエンドポイントの追加

#### Step 1: スキーマ定義

**`api/components/schemas/product.yml`** (新規作成)
```yaml
Product:
  type: object
  description: 商品情報
  properties:
    id:
      type: integer
      format: int64
      description: 商品ID
      example: 1
    name:
      type: string
      description: 商品名
      example: "MacBook Pro"
    price:
      type: integer
      description: 価格（円）
      example: 298000
    categoryId:
      type: integer
      description: カテゴリID
      example: 1
    createdAt:
      type: string
      format: date-time
      description: 作成日時
      example: "2023-12-01T10:00:00Z"

CreateProductRequest:
  type: object
  description: 商品作成リクエスト
  required:
    - name
    - price
    - categoryId
  properties:
    name:
      type: string
      minLength: 1
      description: 商品名
      example: "MacBook Pro"
    price:
      type: integer
      minimum: 0
      description: 価格（円）
      example: 298000
    categoryId:
      type: integer
      description: カテゴリID
      example: 1
```

#### Step 2: パス定義

**`api/paths/products.yml`** (新規作成)
```yaml
products:
  get:
    tags:
      - products
    summary: 商品一覧取得
    description: 商品の一覧を取得します
    operationId: getProducts
    parameters:
      - name: page
        in: query
        description: ページ番号
        schema:
          type: integer
          default: 1
          minimum: 1
      - name: limit
        in: query
        description: 取得件数
        schema:
          type: integer
          default: 20
          minimum: 1
          maximum: 100
    responses:
      "200":
        description: 取得成功
        content:
          application/json:
            schema:
              type: object
              properties:
                products:
                  type: array
                  items:
                    $ref: "../components/schemas/product.yml#/Product"

  post:
    tags:
      - products
    summary: 商品作成
    description: 新しい商品を作成します
    operationId: createProduct
    security:
      - bearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: "../components/schemas/product.yml#/CreateProductRequest"
    responses:
      "201":
        description: 作成成功
        content:
          application/json:
            schema:
              $ref: "../components/schemas/product.yml#/Product"
      "400":
        description: バリデーションエラー
        content:
          application/json:
            schema:
              $ref: "../components/schemas/error.yml#/ErrorResponse"
```

#### Step 3: メインファイルに追加

**`api/index.yml`** に追加
```yaml
paths:
  # 既存のパス...
  
  # Products（追加）
  /products:
    $ref: "./paths/products.yml#/products"

components:
  schemas:
    # 既存のスキーマ...
    
    # Product related schemas（追加）
    Product:
      $ref: "./components/schemas/product.yml#/Product"
    CreateProductRequest:
      $ref: "./components/schemas/product.yml#/CreateProductRequest"
```

#### Step 4: マージ・生成・確認

```bash
# 1. 仕様をマージ
make merge-api

# 2. コード生成
make gen-api

# 3. Swagger UIで確認
make dev
# http://localhost:8080/swagger/ でProducts APIを確認
```

---

## コード生成

### ⚙️ 生成プロセス

#### 1. マージ処理
```bash
make merge-api
```
- 分割されたYAMLファイルを統合
- `api/openapi.yml` に出力
- 参照（`$ref`）を解決

#### 2. コード生成
```bash
make gen-api
```
- OpenAPI Generator を使用
- Go と TypeScript のコードを生成

#### 3. 生成される内容

**Backend (Go) - `backend/api/`**
```
backend/api/
├── go/                          # 生成されたGoコード
│   ├── model_user.go           # ユーザーモデル
│   ├── model_auth_response.go  # 認証レスポンス
│   ├── api_auth.go            # 認証API
│   ├── api_users.go           # ユーザーAPI
│   └── routers.go             # ルーター設定
├── README.md                   # 生成コードの説明
└── .openapi-generator/         # 生成設定
```

**Frontend (TypeScript) - `frontend/src/api/`**
```
frontend/src/api/
├── apis/                       # API client
│   ├── auth-api.ts            # 認証API
│   ├── users-api.ts           # ユーザーAPI
│   └── index.ts               # エクスポート
├── models/                     # 型定義
│   ├── user.ts                # ユーザー型
│   ├── auth-response.ts       # 認証レスポンス型
│   └── index.ts               # エクスポート
├── base.ts                     # 基底クラス
├── common.ts                   # 共通型
├── configuration.ts            # 設定
└── index.ts                    # メインエクスポート
```

### 🔍 生成されたコードの確認

#### Go の例
```go
// backend/api/go/model_user.go
// 自動生成されたコード
type User struct {
    Id        int64     `json:"id"`         
    Name      string    `json:"name"`       
    Email     string    `json:"email"`      
    CreatedAt time.Time `json:"createdAt"`  
    UpdatedAt time.Time `json:"updatedAt"`  
}
```

#### TypeScript の例
```typescript
// frontend/src/api/models/user.ts
// 自動生成されたコード
export interface User {
    id: number;
    name: string;
    email: string;
    createdAt: string;
    updatedAt: string;
}
```

---

## 開発フロー

### 🔄 API-First 開発フロー

#### 1. API設計 → 2. コード生成 → 3. 実装

```
要件定義 → API仕様作成 → make merge-api → Swagger UI確認
    ↓                                         ↓
実装完了 ← Backend実装 ← make gen-api ← 仕様OK?
    ↓        ↓
    ← Frontend実装
```

#### Step 1: API仕様作成
```bash
# 新機能の仕様作成
vim api/paths/products.yml
vim api/components/schemas/product.yml

# 仕様マージ・確認
make merge-api
make dev
# http://localhost:8080/swagger/ で確認
```

#### Step 2: チームレビュー
- Swagger UI でAPIレビュー
- 必要に応じて仕様修正
- レビュー完了後にコード生成

#### Step 3: コード生成・実装
```bash
# コード生成
make gen-api

# Backend実装
vim backend/internal/controller/product_controller.go
vim backend/internal/usecase/product_usecase.go
vim backend/internal/repository/product_repository.go

# Frontend実装
vim frontend/pages/products.tsx
vim frontend/components/ProductList.tsx
```

---

## 実践例

### 💼 実際のAPI追加例：商品管理機能

#### 要件
- 商品の一覧表示・詳細表示・作成・更新・削除
- カテゴリによるフィルタリング
- ページネーション対応
- 管理者のみ作成・更新・削除可能

#### 1. スキーマ設計

**`api/components/schemas/product.yml`**
```yaml
Product:
  type: object
  required:
    - id
    - name
    - price
    - categoryId
  properties:
    id:
      type: integer
      format: int64
      description: 商品ID
    name:
      type: string
      maxLength: 100
      description: 商品名
    description:
      type: string
      maxLength: 1000  
      description: 商品説明
    price:
      type: integer
      minimum: 0
      description: 価格（円）
    categoryId:
      type: integer
      description: カテゴリID
    imageUrl:
      type: string
      format: uri
      description: 商品画像URL
    stock:
      type: integer
      minimum: 0
      description: 在庫数
    isActive:
      type: boolean
      default: true
      description: 販売中かどうか
    createdAt:
      type: string
      format: date-time
      description: 作成日時
    updatedAt:
      type: string
      format: date-time
      description: 更新日時

CreateProductRequest:
  type: object
  required:
    - name
    - price
    - categoryId
  properties:
    name:
      type: string
      minLength: 1
      maxLength: 100
    description:
      type: string
      maxLength: 1000
    price:
      type: integer
      minimum: 0
    categoryId:
      type: integer
    imageUrl:
      type: string
      format: uri
    stock:
      type: integer
      minimum: 0
      default: 0
```

#### 2. 生成・実装

```bash
# コード生成
make gen-api

# Backend実装
vim backend/internal/controller/product_controller.go
vim backend/internal/usecase/product_usecase.go  
vim backend/internal/repository/product_repository.go

# Frontend実装  
vim frontend/pages/products/index.tsx
vim frontend/pages/products/[id].tsx
vim frontend/components/ProductCard.tsx
```

---

## トラブルシューティング

### 🚨 よくある問題と解決法

#### 1. マージエラー

**エラー例:**
```
Error: Cannot resolve reference: ./components/schemas/product.yml#/Product
```

**原因:** 参照パスが間違っている

**解決法:**
```yaml
# ❌ 間違い
$ref: "./schemas/product.yml#/Product"

# ✅ 正しい
$ref: "./components/schemas/product.yml#/Product"
```

#### 2. 生成コードエラー

**エラー例:**
```
Error: Invalid OpenAPI specification
```

**解決法:**
```bash
# 仕様の検証
cd api && npm run validate

# エラー内容を確認して修正
vim api/openapi.yml
```

#### 3. 型エラー

**エラー例:** TypeScriptで生成された型にエラー

**原因:** OpenAPI仕様の型定義が不正確

**解決法:**
```yaml
# ❌ 型が曖昧
price: 
  type: number

# ✅ 明確な型定義
price:
  type: integer
  minimum: 0
  description: 価格（円）
  example: 1000
```

### 🔧 デバッグ方法

#### 1. 段階的確認
```bash
# 1. 個別ファイルの文法チェック
yamllint api/paths/products.yml

# 2. マージ確認
make merge-api
cat api/openapi.yml | head -50

# 3. 仕様検証
cd api && npm run validate
```

#### 2. Swagger UI での確認
```bash
make dev
# http://localhost:8080/swagger/ で仕様を視覚的に確認
```

---

## ベストプラクティス

### ✅ OpenAPI仕様記述のベストプラクティス

#### 1. 詳細な記述

```yaml
# ✅ 推奨: 詳細な記述
properties:
  price:
    type: integer
    minimum: 0
    maximum: 10000000
    description: 商品価格（税込み、円）
    example: 2980

# ❌ 非推奨: 最小限の記述
properties:
  price:
    type: integer
```

#### 2. エラーレスポンスの統一

```yaml
# 全エンドポイントで統一のエラー形式
responses:
  "400":
    description: バリデーションエラー
    content:
      application/json:
        schema:
          $ref: "../components/schemas/error.yml#/ErrorResponse"
```

#### 3. セキュリティの明示

```yaml
# 認証が必要なエンドポイントは明示的に設定
post:
  security:
    - bearerAuth: []
```

### 🏗️ ファイル構成のベストプラクティス

#### 1. 機能別ファイル分割

```
paths/
├── auth.yml          # 認証関連
├── users.yml         # ユーザー管理
├── products.yml      # 商品管理
├── orders.yml        # 注文管理
└── admin.yml         # 管理機能
```

#### 2. スキーマの再利用

```yaml
# ✅ 推奨: 共通スキーマの再利用
components:
  schemas:
    BaseEntity:
      type: object
      properties:
        id:
          type: integer
          format: int64
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
    
    Product:
      allOf:
        - $ref: "#/BaseEntity"
        - type: object
          properties:
            name:
              type: string
            price:
              type: integer
```

### 🔄 開発フローのベストプラクティス

#### 1. API設計レビュー

```bash
# 1. API仕様作成
vim api/paths/new-feature.yml

# 2. マージ・確認
make merge-api
make dev

# 3. チームレビュー（Swagger UI使用）
# 4. 修正・再レビュー
# 5. 承認後にコード生成
make gen-api
```

#### 2. 段階的リリース

```bash
# 1. API仕様のみコミット（実装前）
git add api/
git commit -m "Add product management API specification"

# 2. 生成されたコードをコミット
make gen-api
git add backend/api/ frontend/src/api/
git commit -m "Generate API code from specification"

# 3. 実装後にコミット
git add backend/internal/ frontend/pages/
git commit -m "Implement product management feature"
```

### 🎯 型安全性の確保

#### 1. 厳密な型定義

```yaml
# ✅ 推奨: 厳密な型
properties:
  status:
    type: string
    enum: [active, inactive, pending]
  email:
    type: string
    format: email

# ❌ 非推奨: 曖昧な型
properties:
  status:
    type: string
  email:
    type: string
```

#### 2. 必須フィールドの明示

```yaml
# ✅ 推奨: required を明示
CreateProductRequest:
  type: object
  required:
    - name
    - price
    - categoryId
  properties:
    # ...
```

---

このガイドを参考にして、効率的なAPI設計・開発を行ってください！

質問やトラブルがある場合は、GitHub Issues または Slack で相談してください。 