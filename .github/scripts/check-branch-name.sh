#!/bin/bash

# ブランチ名チェックスクリプト
# Usage: ./check-branch-name.sh <branch-name> <target-branch>

set -e

BRANCH_NAME=$1
TARGET_BRANCH=$2

if [ -z "$BRANCH_NAME" ] || [ -z "$TARGET_BRANCH" ]; then
    echo "❌ Usage: $0 <branch-name> <target-branch>"
    exit 1
fi

echo "🔍 Checking branch name: $BRANCH_NAME → $TARGET_BRANCH"

# mainブランチへのマージルールチェック
if [ "$TARGET_BRANCH" = "main" ]; then
    if [[ ! "$BRANCH_NAME" =~ ^develop/.+ ]]; then
        echo "❌ Error: Pull requests to main branch must come from 'develop/' branches"
        echo "   Current branch: $BRANCH_NAME"
        echo "   Expected pattern: develop/feature-name"
        echo ""
        echo "📋 Examples of valid branch names:"
        echo "   ✅ develop/user-authentication"
        echo "   ✅ develop/fix-login-bug"
        echo "   ✅ develop/add-payment-system"
        echo ""
        echo "📋 Examples of invalid branch names:"
        echo "   ❌ feature/user-auth"
        echo "   ❌ bugfix/login-issue"
        echo "   ❌ user-authentication"
        exit 1
    fi
fi

# developブランチへのマージルールチェック
if [ "$TARGET_BRANCH" = "develop" ]; then
    if [[ "$BRANCH_NAME" =~ ^(feature|fix|hotfix)/.+ ]]; then
        echo "✅ Valid branch name for develop: $BRANCH_NAME"
    else
        echo "⚠️  Warning: Unusual branch name for develop: $BRANCH_NAME"
        echo "   Recommended patterns: feature/*, fix/*, hotfix/*"
    fi
fi

# ブランチ名の一般的なチェック
if [[ "$BRANCH_NAME" =~ [A-Z] ]]; then
    echo "⚠️  Warning: Branch name contains uppercase letters"
    echo "   Consider using lowercase with hyphens: ${BRANCH_NAME,,}"
fi

if [[ "$BRANCH_NAME" =~ _+ ]]; then
    echo "⚠️  Warning: Branch name contains underscores"
    echo "   Consider using hyphens instead of underscores"
fi

echo "✅ Branch name check completed: $BRANCH_NAME" 