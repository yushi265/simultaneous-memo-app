# 開発チェックリスト

新機能追加時や変更時に確認すべき項目のチェックリストです。

## 📝 ドキュメント更新チェックリスト

### 新機能追加時
- [ ] **README.md** - 主な機能セクションに追加
- [ ] **README.md** - 使い方セクションに手順追加
- [ ] **README.md** - API エンドポイントセクションに新しいエンドポイント追加
- [ ] **README.md** - 技術スタックに新しいライブラリ/技術追加
- [ ] **CLAUDE.md** - 実装済み機能リストに追加
- [ ] **package.json** - 新しい依存関係が追加された場合

### UI/コンポーネント変更時
- [ ] **README.md** - プロジェクト構造セクション更新
- [ ] **README.md** - 使い方セクションのスクリーンショット更新（必要に応じて）
- [ ] **CLAUDE.md** - Key Features セクション更新

### API変更時
- [ ] **README.md** - API エンドポイントセクション更新
- [ ] **CLAUDE.md** - API エンドポイントセクション更新
- [ ] Postmanコレクションまたはswagger.yml更新（あれば）

### 設定変更時
- [ ] **README.md** - 起動方法・設定方法更新
- [ ] **docker-compose.yml** - 環境変数追加時はREADMEにも記載
- [ ] **.env.example** - 新しい環境変数追加

## 🧪 テスト関連チェックリスト

### 新機能追加時
- [ ] 単体テスト追加
- [ ] 統合テスト追加（API変更時）
- [ ] E2Eテスト追加（重要な機能の場合）
- [ ] 手動テストの実施

### 既存機能変更時
- [ ] 関連テストの更新
- [ ] 回帰テストの実施
- [ ] パフォーマンステスト（必要に応じて）

## 🔄 PR作成前チェックリスト

### コード品質
- [ ] Lintエラーがないことを確認
- [ ] Type errorがないことを確認
- [ ] ビルドが成功することを確認
- [ ] 不要なコンソールログ・デバッグコードの削除

### ドキュメント
- [ ] 上記のドキュメント更新チェックリストを実行
- [ ] コミットメッセージが適切
- [ ] PR説明が詳細

### テスト
- [ ] 新機能が期待通り動作する
- [ ] 既存機能に影響がない
- [ ] エラーハンドリングが適切

## 📋 定期メンテナンスチェックリスト（月次）

### 依存関係
- [ ] npm audit で脆弱性チェック
- [ ] go mod で依存関係更新チェック
- [ ] 未使用の依存関係削除

### ドキュメント
- [ ] README.mdの情報が最新か確認
- [ ] リンクが有効か確認
- [ ] スクリーンショットが最新か確認

### パフォーマンス
- [ ] バンドルサイズチェック
- [ ] ページ読み込み速度チェック
- [ ] データベースパフォーマンス確認

## 💡 便利なコマンド

```bash
# ドキュメント更新忘れチェック
git diff --name-only HEAD~1 | grep -E "(package\.json|go\.mod|main\.go)" && echo "Consider updating README.md"

# 新しいコンポーネント追加チェック  
git diff --name-only HEAD~1 | grep "frontend/components/" && echo "Update README project structure"

# API変更チェック
git diff HEAD~1 -- backend/main.go | grep "api\." && echo "Update README API endpoints"

# ドキュメントの整合性チェック
grep -o "frontend/components/[^)]*" README.md | while read file; do
  [ ! -f "$file" ] && echo "Missing file: $file"
done
```

## 🔗 関連リソース

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)