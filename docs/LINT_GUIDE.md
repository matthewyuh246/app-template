# 📋 Lint設定ガイド

このドキュメントでは、プロジェクトで使用しているLint（コード品質チェック）ツールの設定と使い方について説明します。

## 📖 目次

- [🎯 Lintとは](#-lintとは)
- [🔧 バックエンド（Go）Lint](#-バックエンドgolint)
- [🎨 フロントエンド（TypeScript/React）Lint](#-フロントエンドtypescriptreactlint)
- [🚀 CI/CDでのLint実行](#-cicdでのlint実行)
- [💡 よくあるエラーと対処法](#-よくあるエラーと対処法)
- [⚙️ カスタマイズ方法](#️-カスタマイズ方法)

---

## 🎯 Lintとは

**Lint（リント）** は、コードの品質を自動的にチェックするツールです。

### 🔍 Lintが検出する問題

| カテゴリ | 説明 | 例 |
|---------|------|-----|
| **構文エラー** | コードの文法的な間違い | セミコロン忘れ、括弧の不一致 |
| **スタイル** | コードの書き方の統一 | インデント、命名規則 |
| **潜在的バグ** | 実行時エラーの可能性 | 未使用変数、型の不一致 |
| **セキュリティ** | セキュリティ上の問題 | SQLインジェクション、XSS |
| **パフォーマンス** | 性能上の問題 | 非効率なループ、メモリリーク |

### ✅ Lintのメリット

- **🐛 バグの早期発見**: 実行前に問題を検出
- **📏 コード品質の統一**: チーム全体で一貫したコードスタイル
- **🚀 開発効率の向上**: 自動チェックで手動レビューの負荷軽減
- **📚 学習効果**: ベストプラクティスの習得

---

## 🔧 バックエンド（Go）Lint

### 🛠️ 使用ツール: golangci-lint

**golangci-lint** は、複数のGoのLintツールを統合した高性能なLintツールです。

#### 📁 設定ファイル: `backend/.golangci.yml`

```yaml
# Lintの実行設定
run:
  timeout: 5m                    # タイムアウト時間
  tests: true                    # テストファイルもチェック
  skip-dirs:                     # チェック除外ディレクトリ
    - vendor
    - tmp

# 有効にするLinter
linters:
  enable:
    - gofmt          # コードフォーマット
    - goimports      # importの整理
    - govet          # 潜在的バグ検出
    - ineffassign    # 無効な代入検出
    - misspell       # スペルミス検出
    - unparam        # 未使用パラメータ検出
    - unused         # 未使用コード検出
    - errcheck       # エラーハンドリング検証
    - gosimple       # コード簡略化提案
    - staticcheck    # 静的解析
    - typecheck      # 型チェック
    - gocritic       # パフォーマンス・スタイル
    - revive         # 高速なgo-lint代替

# Linterの詳細設定
linters-settings:
  gofmt:
    simplify: true             # コード簡略化
  goimports:
    local-prefixes: github.com/yourorg/app-template
  revive:
    min-confidence: 0.8        # 信頼度の閾値
  gocritic:
    enabled-tags:
      - diagnostic             # 診断系
      - style                  # スタイル系
      - performance            # パフォーマンス系
```

#### 🚀 実行方法

**ローカル実行:**
```bash
# backend ディレクトリで実行
cd backend

# 全ファイルをチェック
golangci-lint run

# 特定のディレクトリをチェック
golangci-lint run ./internal/...

# 自動修正可能な問題を修正
golangci-lint run --fix
```

**CI実行:**
```bash
# CIで自動実行（修正なし）
golangci-lint run --timeout=5m
```

#### 📋 主要なチェック項目

##### 1️⃣ **コードフォーマット（gofmt）**
```go
// ❌ 悪い例
func   BadFunction(  x int,y int  )int{
return x+y
}

// ✅ 良い例
func GoodFunction(x int, y int) int {
    return x + y
}
```

##### 2️⃣ **Import整理（goimports）**
```go
// ❌ 悪い例
import (
    "github.com/gin-gonic/gin"
    "fmt"
    "github.com/yourorg/app-template/internal/entity"
    "database/sql"
)

// ✅ 良い例
import (
    "database/sql"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/yourorg/app-template/internal/entity"
)
```

##### 3️⃣ **エラーハンドリング（errcheck）**
```go
// ❌ 悪い例
func BadExample() {
    file, err := os.Open("file.txt")
    file.Close() // エラーチェックなし
}

// ✅ 良い例
func GoodExample() error {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := file.Close(); closeErr != nil {
            log.Printf("Failed to close file: %v", closeErr)
        }
    }()
    return nil
}
```

##### 4️⃣ **未使用変数/関数（unused）**
```go
// ❌ 悪い例
func BadExample() {
    unusedVar := "never used"     // 未使用変数
    x := calculateValue()
    y := anotherValue()           // 未使用変数
    return x
}

// ✅ 良い例
func GoodExample() {
    x := calculateValue()
    return x
}
```

##### 5️⃣ **命名規則（revive）**
```go
// ❌ 悪い例
func get_user_data() {}           // スネークケース
var HTTPSserver *http.Server      // 略語の大文字化
const max_retry_count = 3         // 定数のスネークケース

// ✅ 良い例
func GetUserData() {}             // キャメルケース
var HTTPSServer *http.Server      // 略語は全て大文字
const MaxRetryCount = 3           // 定数もキャメルケース
```

---

## 🎨 フロントエンド（TypeScript/React）Lint

### 🛠️ 使用ツール: ESLint + TypeScript

**ESLint** は、JavaScriptとTypeScriptのためのLintツールです。

#### 📁 設定ファイル: `frontend/.eslintrc.json`

```json
{
  "extends": [
    "next/core-web-vitals",      // Next.js推奨設定
    "@typescript-eslint/recommended", // TypeScript推奨設定
    "plugin:react/recommended",   // React推奨設定
    "plugin:react-hooks/recommended" // React Hooks推奨設定
  ],
  "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaVersion": 2022,
    "sourceType": "module",
    "ecmaFeatures": {
      "jsx": true
    }
  },
  "plugins": [
    "@typescript-eslint",
    "react",
    "react-hooks"
  ],
  "rules": {
    // TypeScript関連
    "@typescript-eslint/no-unused-vars": "error",
    "@typescript-eslint/no-explicit-any": "warn",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/no-non-null-assertion": "warn",
    
    // React関連
    "react/react-in-jsx-scope": "off",
    "react/prop-types": "off",
    "react-hooks/rules-of-hooks": "error",
    "react-hooks/exhaustive-deps": "warn",
    
    // 一般的なJavaScript/TypeScript
    "no-console": "warn",
    "no-debugger": "error",
    "prefer-const": "error",
    "no-var": "error"
  },
  "settings": {
    "react": {
      "version": "detect"
    }
  }
}
```

#### 🚀 実行方法

**ローカル実行:**
```bash
# frontend ディレクトリで実行
cd frontend

# 全ファイルをチェック
npm run lint

# 自動修正可能な問題を修正
npm run lint -- --fix

# 特定のファイルをチェック
npx eslint pages/index.tsx

# TypeScript型チェック
npm run type-check
```

**CI実行:**
```bash
# CIで自動実行
npm run lint
npm run type-check
```

#### 📋 主要なチェック項目

##### 1️⃣ **TypeScript型安全性**
```typescript
// ❌ 悪い例
function badFunction(data: any) {          // any型の使用
    return data.someProperty;
}

const user = getUserData();
console.log(user.name);                    // 型チェックなし

// ✅ 良い例
interface User {
    id: number;
    name: string;
    email: string;
}

function goodFunction(data: User): string {
    return data.name;
}

const user: User = getUserData();
console.log(user.name);                    // 型安全
```

##### 2️⃣ **React Hooks ルール**
```typescript
// ❌ 悪い例
function BadComponent({ userId }: { userId: string }) {
    if (userId) {
        const [user, setUser] = useState<User | null>(null); // 条件内でHook使用
    }
    
    useEffect(() => {
        fetchUser();
    }, []); // 依存配列が不完全
    
    return <div>{user?.name}</div>;
}

// ✅ 良い例
function GoodComponent({ userId }: { userId: string }) {
    const [user, setUser] = useState<User | null>(null);
    
    useEffect(() => {
        if (userId) {
            fetchUser(userId);
        }
    }, [userId]); // 正しい依存配列
    
    return <div>{user?.name}</div>;
}
```

##### 3️⃣ **未使用変数・Import**
```typescript
// ❌ 悪い例
import React, { useState, useEffect } from 'react';
import { Button } from '../components/Button';    // 未使用import
import axios from 'axios';

function BadComponent() {
    const [data, setData] = useState(null);
    const unusedVariable = 'never used';          // 未使用変数
    
    return <div>Hello World</div>;
}

// ✅ 良い例
import React, { useState } from 'react';

function GoodComponent() {
    const [data, setData] = useState(null);
    
    return <div>Hello World</div>;
}
```

##### 4️⃣ **コンソール・デバッガー**
```typescript
// ❌ 悪い例（本番環境で残ってはいけない）
function badFunction() {
    console.log('Debug info');                     // 警告
    debugger;                                      // エラー
    
    return processData();
}

// ✅ 良い例
function goodFunction() {
    // 本番環境では適切なログライブラリを使用
    return processData();
}
```

##### 5️⃣ **現代的なJavaScript**
```typescript
// ❌ 悪い例
var oldStyleVariable = 'bad';                      // var使用
let mutableButNeverChanged = 'should be const';    // let不適切使用

function oldFunction() {
    var result = [];
    for (var i = 0; i < items.length; i++) {
        result.push(items[i].name);
    }
    return result;
}

// ✅ 良い例
const modernVariable = 'good';                     // const使用

function modernFunction() {
    return items.map(item => item.name);           // 関数型プログラミング
}
```

---

## 🚀 CI/CDでのLint実行

### 🔄 自動実行（PR作成時）

#### バックエンドLint
```yaml
- name: Run linting
  working-directory: ./backend
  run: |
    echo "🔍 Running Go linting..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    golangci-lint run --timeout=5m
    echo "✅ Linting passed"
```

#### フロントエンドLint
```yaml
- name: Run ESLint
  working-directory: ./frontend
  run: |
    echo "🔍 Running ESLint..."
    npm run lint
    echo "✅ ESLint passed"

- name: Run type check
  working-directory: ./frontend
  run: |
    echo "🔍 Running TypeScript type check..."
    npm run type-check
    echo "✅ Type check passed"
```

### 📊 CI実行結果の見方

#### ✅ 成功例
```
🔍 Running Go linting...
✅ Linting passed

🔍 Running ESLint...
✅ ESLint passed

🔍 Running TypeScript type check...
✅ Type check passed
```

#### ❌ 失敗例
```
🔍 Running Go linting...
backend/internal/controller/user_controller.go:15:2: 
  unused variable 'unusedVar' (unused)
backend/internal/entity/user.go:25:1: 
  exported function GetUser should have comment (revive)

❌ Linting failed
```

---

## 💡 よくあるエラーと対処法

### 🔧 バックエンド（Go）

#### 1️⃣ **未使用変数エラー**
```
Error: unused variable 'result' (unused)
```

**対処法:**
```go
// 方法1: 変数を削除
func fixedFunction() {
    // result := calculateValue() // 削除
    return
}

// 方法2: アンダースコアで無視
func fixedFunction() {
    _ = calculateValue() // 意図的に無視
    return
}

// 方法3: 実際に使用
func fixedFunction() {
    result := calculateValue()
    log.Printf("Result: %v", result) // 使用
    return
}
```

#### 2️⃣ **エクスポート関数のコメントなし**
```
Error: exported function GetUser should have comment (revive)
```

**対処法:**
```go
// ❌ 悪い例
func GetUser(id string) (*User, error) {
    // ...
}

// ✅ 良い例
// GetUser retrieves a user by ID from the database
func GetUser(id string) (*User, error) {
    // ...
}
```

#### 3️⃣ **エラーハンドリング不備**
```
Error: Error return value of 'file.Close' is not checked (errcheck)
```

**対処法:**
```go
// ❌ 悪い例
file, err := os.Open("file.txt")
if err != nil {
    return err
}
file.Close() // エラーチェックなし

// ✅ 良い例
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer func() {
    if err := file.Close(); err != nil {
        log.Printf("Failed to close file: %v", err)
    }
}()
```

### 🎨 フロントエンド（TypeScript/React）

#### 1️⃣ **型エラー**
```
Error: Property 'name' does not exist on type 'unknown'
```

**対処法:**
```typescript
// ❌ 悪い例
function badFunction(data: unknown) {
    return data.name; // エラー: unknownに'name'プロパティは存在しない
}

// ✅ 良い例
interface User {
    name: string;
}

function goodFunction(data: User) {
    return data.name; // OK: 型が明確
}

// または型ガードを使用
function safeFunction(data: unknown) {
    if (typeof data === 'object' && data !== null && 'name' in data) {
        return (data as { name: string }).name;
    }
    return null;
}
```

#### 2️⃣ **React Hooks依存配列**
```
Warning: React Hook useEffect has a missing dependency: 'userId'
```

**対処法:**
```typescript
// ❌ 悪い例
function BadComponent({ userId }: { userId: string }) {
    const [user, setUser] = useState(null);
    
    useEffect(() => {
        fetchUser(userId);
    }, []); // userIdが依存配列にない
    
    return <div>{user?.name}</div>;
}

// ✅ 良い例
function GoodComponent({ userId }: { userId: string }) {
    const [user, setUser] = useState(null);
    
    useEffect(() => {
        fetchUser(userId);
    }, [userId]); // 正しい依存配列
    
    return <div>{user?.name}</div>;
}
```

#### 3️⃣ **未使用Import**
```
Error: 'Button' is defined but never used (no-unused-vars)
```

**対処法:**
```typescript
// ❌ 悪い例
import React from 'react';
import { Button } from '../components/Button'; // 未使用

function Component() {
    return <div>Hello</div>;
}

// ✅ 良い例（方法1: Importを削除）
import React from 'react';

function Component() {
    return <div>Hello</div>;
}

// ✅ 良い例（方法2: 実際に使用）
import React from 'react';
import { Button } from '../components/Button';

function Component() {
    return (
        <div>
            <Button>Click me</Button>
        </div>
    );
}
```

---

## ⚙️ カスタマイズ方法

### 🔧 バックエンドLint設定のカスタマイズ

#### 特定のルールを無効化
```yaml
# backend/.golangci.yml
linters-settings:
  revive:
    rules:
      - name: exported-function-comment
        disabled: true  # エクスポート関数のコメント必須を無効化
```

#### ファイルを除外
```yaml
# backend/.golangci.yml
issues:
  exclude-rules:
    - path: ".*_test.go"
      linters:
        - errcheck  # テストファイルではerrcheckを無効化
```

#### 特定の行を除外
```go
//nolint:errcheck // この行のerrcheckを無視
file.Close()

//nolint:all // この行の全てのLintを無視
riskyOperation()
```

### 🎨 フロントエンドLint設定のカスタマイズ

#### 特定のルールを無効化
```json
{
  "rules": {
    "no-console": "off",  // console.logを許可
    "@typescript-eslint/no-explicit-any": "off"  // any型を許可
  }
}
```

#### ファイルを除外
```json
{
  "ignorePatterns": [
    "dist/",
    "build/",
    "*.config.js"
  ]
}
```

#### 特定の行を除外
```typescript
// eslint-disable-next-line no-console
console.log('Debug information');

/* eslint-disable @typescript-eslint/no-explicit-any */
function legacyFunction(data: any) {
    return data;
}
/* eslint-enable @typescript-eslint/no-explicit-any */
```

---

## 🎯 ベストプラクティス

### ✅ 推奨事項

1. **🔄 定期的なLint実行**
   - コミット前に必ずローカルでLint実行
   - IDE/エディタのLint拡張機能を活用

2. **📏 一貫したルール適用**
   - チーム全体で同じLint設定を使用
   - 設定変更は全員で議論して決定

3. **🎓 学習機会として活用**
   - Lintエラーから新しいベストプラクティスを学ぶ
   - 警告も無視せず、理由を理解して対処

4. **⚡ 効率的な修正**
   - 自動修正機能を積極的に活用
   - 似たようなエラーはパターンを覚えて予防

### ❌ 避けるべき事項

1. **🚫 Lintルールの無効化乱用**
   - 安易に`//nolint`や`eslint-disable`を使わない
   - 無効化する場合は理由をコメントで明記

2. **🚫 CI失敗の放置**
   - Lint失敗は必ず修正してからマージ
   - 「後で直す」は禁物

3. **🚫 設定の個人的変更**
   - 個人の好みで勝手にルールを変更しない
   - 変更が必要な場合はチームで相談

---

## 📚 参考リンク

### 🔧 バックエンド（Go）
- [golangci-lint公式ドキュメント](https://golangci-lint.run/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 🎨 フロントエンド（TypeScript/React）
- [ESLint公式ドキュメント](https://eslint.org/docs/user-guide/)
- [TypeScript ESLint](https://typescript-eslint.io/)
- [React ESLint Plugin](https://github.com/yannickcr/eslint-plugin-react)

---

**最終更新**: 2024年12月
**バージョン**: v1.0
**対象**: Go 1.24, TypeScript 5.x, React 18.x 