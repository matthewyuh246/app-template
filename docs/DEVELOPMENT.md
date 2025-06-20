# 開発ガイド

## 🌳 ブランチ戦略

このプロジェクトでは **Git Flow** を簡略化したブランチ戦略を採用しています。

### ブランチ構成

```
main (本番環境)
 ↑
 └─ develop/feature-name (開発ブランチ)
     ↑
     ├─ feature/user-auth (機能開発)
     ├─ fix/login-bug (バグ修正)
     └─ hotfix/security-patch (緊急修正)
```

### ブランチルール

#### 1. **main ブランチ**
- 本番環境のコード
- **直接pushは禁止**
- `develop/` から始まるブランチからのみマージ可能
- 全てのCIチェックが必須
- コードオーナーのレビュー必須

#### 2. **develop/\* ブランチ**
- 開発作業用のブランチ
- `develop/feature-name` の形式で命名
- mainブランチへのPRに使用

#### 3. **feature/\* ブランチ**
- 新機能開発用
- `feature/feature-name` の形式で命名
- developブランチにマージ

#### 4. **fix/\* ブランチ**
- バグ修正用
- `fix/bug-description` の形式で命名

#### 5. **hotfix/\* ブランチ**
- 緊急修正用
- `hotfix/issue-description` の形式で命名

## 🔄 開発フロー

### 1. 新機能開発

```bash
# 1. develop系ブランチを作成
git checkout main
git pull origin main
git checkout -b develop/user-authentication

# 2. 機能を開発
# ... コード作成 ...

# 3. テストとlint実行
make test
make lint

# 4. コミットとプッシュ
git add .
git commit -m "feat: ユーザー認証機能を追加"
git push origin develop/user-authentication

# 5. Pull Request作成（main ← develop/user-authentication）
```

### 2. バグ修正

```bash
# 1. fix系ブランチを作成
git checkout main
git pull origin main
git checkout -b fix/login-redirect-issue

# 2. バグを修正
# ... コード修正 ...

# 3. テストとlint実行
make test
make lint

# 4. develop系ブランチとしてPR作成
git checkout -b develop/fix-login-redirect-issue
git cherry-pick <fix-commit-hash>
git push origin develop/fix-login-redirect-issue
```

## 🚫 NG パターン

### ❌ mainブランチへの直接PR
```bash
# これはNG - CIで失敗します
feature/user-auth → main
fix/bug-123 → main
user-auth → main
```

### ✅ OK パターン
```bash
# これがOK - CIが通ります
develop/user-auth → main
develop/fix-bug-123 → main
develop/add-payment-system → main
```

## 📋 Pull Request ガイドライン

### PR作成前チェックリスト

- [ ] ブランチ名が `develop/` で始まっている
- [ ] 全てのテストが通っている
- [ ] Lintエラーがない
- [ ] コンフリクトが解決されている
- [ ] 適切なコミットメッセージ

### PRタイトル形式

```
feat: 新機能の追加
fix: バグの修正
docs: ドキュメントの更新
style: コードフォーマット
refactor: リファクタリング
test: テストの追加・修正
chore: その他の変更
```

### PR説明文

PRテンプレートを使用して以下を記載：

- 変更内容の概要
- テスト内容
- スクリーンショット（UI変更の場合）
- レビューポイント

## 🔄 CI/CD パイプライン

### Pull Request時

1. **ブランチ名チェック**
   - mainへのPRは `develop/` ブランチのみ許可

2. **変更検出**
   - 変更されたファイルに応じてテストを実行

3. **テスト実行**
   - バックエンド: Go テスト + Lint
   - フロントエンド: Jest + ESLint + Type Check

4. **セキュリティスキャン**
   - 脆弱性チェック
   - シークレットスキャン

5. **統合テスト**
   - Docker Compose で統合テスト

6. **PR要約コメント**
   - 結果をPRにコメント

### mainブランチマージ時

1. **最終チェック**
   - develop系ブランチからのマージか確認

2. **ビルド**
   - 本番用ビルド作成

3. **セキュリティチェック**
   - 依存関係脆弱性スキャン

4. **デプロイ**
   - AWS ECS/EKSへの自動デプロイ

5. **スモークテスト**
   - 本番環境での基本動作確認

## 🛠️ 開発環境セットアップ

### 必須ツール

- Docker & Docker Compose
- Git
- Make

### セットアップ手順

```bash
# 1. リポジトリクローン
git clone <repository-url>
cd App-Template

# 2. 環境変数設定
cp env.example .env

# 3. 開発環境起動
make setup
make dev

# 4. 動作確認
curl http://localhost:8080/health
curl http://localhost:3000
```

## 🧪 テスト戦略

### バックエンド

```bash
# 単体テスト
cd backend && go test ./...

# カバレッジ
cd backend && go test -cover ./...

# Lint
cd backend && golangci-lint run
```

### フロントエンド

```bash
# 単体テスト
cd frontend && npm test

# 型チェック
cd frontend && npm run type-check

# Lint
cd frontend && npm run lint

# ビルド
cd frontend && npm run build
```

### 統合テスト

```bash
# 統合テスト実行
make test-integration
```

## 📊 コード品質

### 品質指標

- **テストカバレッジ**: 80%以上
- **Lint**: エラー0個
- **型安全性**: TypeScript strict mode
- **セキュリティ**: 脆弱性0個

### コードレビューガイドライン

#### レビュー観点

1. **機能性**: 要件を満たしているか
2. **可読性**: コードが理解しやすいか
3. **保守性**: 変更・拡張しやすいか
4. **パフォーマンス**: 性能に問題がないか
5. **セキュリティ**: セキュリティホールがないか

#### レビュー時間

- **小さな変更**: 24時間以内
- **中程度の変更**: 48時間以内
- **大きな変更**: 72時間以内

## 🚨 トラブルシューティング

### CIエラーの対処

#### ブランチ名エラー
```
❌ Error: Pull requests to main branch must come from 'develop/' branches
```

**解決策**:
```bash
# 新しいdevelop系ブランチを作成
git checkout -b develop/your-feature-name
git cherry-pick <your-commits>
git push origin develop/your-feature-name
```

#### テストエラー
```
❌ Test failed: 1 failing
```

**解決策**:
```bash
# ローカルでテスト実行
make test

# 失敗したテストを修正
# 再度プッシュ
```

#### Lintエラー
```
❌ Lint failed: 5 errors
```

**解決策**:
```bash
# 自動修正
make lint-fix

# 手動修正が必要な場合は個別に対応
```

## 🔐 セキュリティ

### 機密情報の管理

- `.env` ファイルは Git にコミットしない
- APIキーは GitHub Secrets で管理
- データベースパスワードは環境変数で管理

### セキュリティチェック

- Dependabot による依存関係更新
- Trivy による脆弱性スキャン
- TruffleHog によるシークレットスキャン

## 📞 サポート

### 質問・相談

- GitHub Issues でバグ報告・機能要望
- GitHub Discussions で質問・相談
- Slack `#development` チャンネル

### ドキュメント

- [API ドキュメント](http://localhost:8081) (Swagger UI)
- [プロジェクト概要](../README.md)
- [デプロイガイド](./DEPLOYMENT.md) 