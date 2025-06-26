# リアルタイムメモ

複数ユーザーが同時にリアルタイムで編集可能なNotionライクなメモアプリケーション

![Logo](frontend/public/logo.svg)

## ✨ 主な機能

- 🔐 **ユーザー認証**: JWT認証によるセキュアなログイン・ユーザー管理
- 🏢 **ワークスペース管理**: 個人・チーム用ワークスペースの作成・切り替え
- 📝 **リッチテキストエディター**: 見出し、リスト、コードブロック、太字・斜体などの豊富なフォーマット
- 🔄 **リアルタイム同期**: 複数ユーザーが同時編集可能（Yjs CRDT使用）
- 👥 **協調カーソル**: 他のユーザーのカーソル位置をリアルタイム表示
- 💾 **自動保存**: 3秒のデバウンスで自動保存（レート制限対応）
- 🖼️ **画像機能**: ドラッグ&ドロップ、クリップボード、ボタンからの画像アップロード
- 🎛️ **画像編集**: エディター内でのリサイズ、レスポンシブ配信、サムネイル自動生成
- 📎 **ファイルアップロード**: PDF、ドキュメント、アーカイブ、コードファイルのアップロード対応
- 📁 **ファイル管理**: アップロードファイルの一覧表示、ダウンロード、削除機能
- 🌐 **日本語対応**: 完全日本語化されたUI
- 🚀 **パフォーマンス**: リクエストキャッシュ・リトライ機能でスムーズな操作

## 🛠 技術スタック

### フロントエンド
- **Next.js 14** (App Router)
- **TypeScript** - 型安全性
- **Tailwind CSS v3** - スタイリング
- **TipTap v2** - リッチテキストエディター（画像拡張含む）
- **Yjs** - リアルタイム同期（CRDT）
- **Zustand** - 状態管理（永続化対応）
- **Radix UI** - アイコン

### バックエンド
- **Go 1.23** + **Echo v4** - API フレームワーク
- **PostgreSQL 16** - データベース（JSONB、UUID使用）
- **GORM v2** - ORM（datatypes、hooks対応）
- **JWT** - 認証・認可システム
- **bcrypt** - パスワードハッシュ化
- **Gorilla WebSocket** - リアルタイム通信
- **disintegration/imaging** - 画像処理・リサイズ
- **Air** - ホットリロード

### インフラ
- **Docker** - 開発環境
- **Docker Compose** - マルチコンテナ管理

## 🚀 クイックスタート

### 前提条件
- Docker & Docker Compose

### 起動方法

```bash
# リポジトリをクローン
git clone https://github.com/yushi265/simultaneous-memo-app.git
cd simultaneous-memo-app

# Dockerコンテナを起動
docker-compose up -d --build

# 起動確認
docker-compose ps
```

### アクセス

- **フロントエンド**: <http://localhost:3000>
- **バックエンドAPI**: <http://localhost:8080>
- **ヘルスチェック**: <http://localhost:8080/health>

## 🎯 使い方

### 基本的な使い方

1. ブラウザで <http://localhost:3000> にアクセス
2. **初回利用時**: 「登録」からアカウントを作成（個人ワークスペースが自動作成されます）
3. **既存ユーザー**: 「ログイン」でサインイン
4. 「新規ページ」ボタンをクリックしてページを作成
5. タイトルと本文を編集（自動保存されます）

### ワークスペース機能

- **ワークスペース切り替え**: ヘッダーのワークスペース名をクリック
- **新規ワークスペース作成**: ワークスペースドロップダウンから「新しいワークスペース」
- **ワークスペース設定**: ユーザーメニュー → 「ワークスペース設定」

### コンテンツ編集

- **画像アップロード**:
  - ツールバーの画像ボタンをクリック
  - エディターに画像をドラッグ&ドロップ
  - クリップボードから画像をペースト（Ctrl+V/Cmd+V）
- **画像編集**: 画像を選択してリサイズハンドルでサイズ調整
- **ファイルアップロード**: ツールバーのファイルボタンまたはドラッグ&ドロップ

### リアルタイム協調編集

- 複数のブラウザタブまたは異なるデバイスで同じページにアクセス
- 他のユーザーのカーソル位置と編集内容がリアルタイムで表示

## 📝 開発者向け情報

### 個別に起動する場合

```bash
# フロントエンド開発サーバー
cd frontend
npm install
npm run dev

# バックエンド開発サーバー
cd backend
go mod download
air

# PostgreSQL（別途起動が必要）
docker run --name postgres -e POSTGRES_PASSWORD=dev123 -e POSTGRES_DB=notion_app -p 5432:5432 -d postgres:16
```

### 便利なコマンド

```bash
# ログ確認
docker-compose logs frontend
docker-compose logs backend

# 特定のサービスを再起動
docker-compose restart frontend

# 開発環境をクリーンアップ
docker-compose down -v

# ドキュメント更新リマインダーの設定（推奨）
./scripts/setup-hooks.sh
```

## 📡 API エンドポイント

### 認証
- `POST /api/auth/register` - ユーザー登録（個人ワークスペース自動作成）
- `POST /api/auth/login` - ログイン
- `POST /api/auth/logout` - ログアウト
- `GET /api/auth/me` - ユーザー情報取得

### ワークスペース管理
- `GET /api/workspaces` - ワークスペース一覧取得
- `POST /api/workspaces` - ワークスペース作成
- `GET /api/workspaces/:id` - ワークスペース詳細取得
- `PUT /api/workspaces/:id` - ワークスペース更新
- `DELETE /api/workspaces/:id` - ワークスペース削除
- `POST /api/workspaces/:id/switch` - ワークスペース切り替え

### ページ管理
- `GET /api/pages` - ページ一覧取得（ワークスペース別）
- `POST /api/pages` - ページ作成
- `GET /api/pages/:id` - ページ詳細取得
- `PUT /api/pages/:id` - ページ更新
- `DELETE /api/pages/:id` - ページ削除

### 画像管理
- `POST /api/upload` - 画像アップロード（ページID関連付け対応）
- `GET /api/img/*` - レスポンシブ画像配信（サムネイル対応）
- `GET /api/images` - 画像一覧取得
- `GET /api/images/:id` - 特定画像の詳細取得
- `DELETE /api/images/:id` - 画像削除
- `POST /api/admin/cleanup-images` - 孤立画像のクリーンアップ

### ファイル管理
- `POST /api/upload/file` - 汎用ファイルアップロード
- `GET /api/files` - ファイル一覧取得（フィルタリング対応）
- `GET /api/files/:id` - ファイルメタデータ取得
- `DELETE /api/files/:id` - ファイル削除
- `GET /api/file/*` - ファイル配信

### リアルタイム通信
- `WebSocket /ws/:pageId` - リアルタイム同期（認証対応）

## 📊 システム設計

- [📋 データベーススキーマ](./docs/database-schema.md) - ER図とテーブル設計
- [🏗 システムアーキテクチャ](./docs/architecture.md) - 全体設計とデータフロー

## 🏗 プロジェクト構造

```
├── frontend/                 # Next.js フロントエンド
│   ├── app/                 # App Router
│   │   ├── login/           # ログインページ
│   │   ├── register/        # 登録ページ
│   │   └── workspace/       # ワークスペース管理
│   ├── components/          # Reactコンポーネント
│   │   ├── Header.tsx       # ヘッダー（ユーザーメニュー付き）
│   │   ├── Sidebar.tsx      # サイドバー
│   │   ├── AuthGuard.tsx    # 認証ガード
│   │   ├── WorkspaceSwitcher.tsx # ワークスペース切り替え
│   │   ├── CreateWorkspaceModal.tsx # ワークスペース作成
│   │   ├── Editor.tsx       # エディター（画像・ファイル対応）
│   │   ├── EditorMenuBar.tsx # エディターツールバー
│   │   ├── ResizableImage.tsx # リサイズ可能画像コンポーネント
│   │   ├── FileUpload.tsx   # ファイルアップロードコンポーネント
│   │   └── Logo.tsx         # ロゴ
│   ├── lib/                 # ユーティリティ
│   │   ├── api.ts           # APIクライアント（リトライ対応）
│   │   ├── auth-api.ts      # 認証APIクライアント
│   │   ├── workspace-api.ts # ワークスペースAPIクライアント
│   │   ├── store.ts         # Zustand状態管理（認証・ワークスペース）
│   │   ├── image-upload.ts  # 画像アップロード処理
│   │   ├── image-utils.ts   # 画像関連ユーティリティ
│   │   └── image-resize-extension.ts # TipTap画像拡張
│   └── public/              # 静的ファイル
├── backend/                 # Go バックエンド
│   ├── config/              # 設定管理
│   ├── models/              # データモデル
│   │   ├── user.go          # ユーザー・ワークスペースモデル
│   │   ├── page.go          # ページモデル
│   │   ├── image.go         # 画像モデル
│   │   └── file.go          # ファイルモデル
│   ├── handlers/            # HTTPハンドラー
│   │   ├── auth.go          # 認証ハンドラー
│   │   ├── workspace.go     # ワークスペースハンドラー
│   │   ├── page.go          # ページハンドラー
│   │   ├── image*.go        # 画像処理関連ハンドラー
│   │   ├── file.go          # 画像アップロード
│   │   └── file_general.go  # 汎用ファイルアップロード
│   ├── middleware/          # ミドルウェア
│   │   └── rate_limit.go    # レート制限
│   └── websocket/           # WebSocket処理
├── uploads/                 # アップロードファイル
│   ├── images/              # 画像ファイル（YYYY/MM構造）
│   └── files/               # 汎用ファイル（YYYY/MM構造）
├── docs/                    # プロジェクトドキュメント
│   ├── database-schema.md   # データベース設計
│   └── architecture.md     # システムアーキテクチャ
└── docker-compose.yml       # Docker設定
```

## 📄 ライセンス

MIT License

## 🖼️ 画像機能の詳細

### サポートされている画像形式
- JPEG/JPG
- PNG
- GIF
- WebP

### 画像機能
- **アップロード方法**: ボタンクリック、ドラッグ&ドロップ、クリップボードペースト
- **自動リサイズ**: アップロード時に最適化とサムネイル生成
- **インタラクティブリサイズ**: エディター内でドラッグハンドルによるサイズ調整
- **レスポンシブ配信**: デバイスに応じた最適なサイズで配信
- **メタデータ管理**: ファイルサイズ、寸法、アップロード日時の自動記録
- **リアルタイム同期**: 画像の追加・編集・削除がリアルタイムで他のユーザーに反映

### 制限事項
- 最大ファイルサイズ: 10MB
- 対応画像形式: JPEG、PNG、GIF、WebP

## 📁 ファイル機能の詳細

### サポートされているファイル形式

**ドキュメント**
- PDF、DOC、DOCX、XLS、XLSX、PPT、PPTX
- TXT、CSV、RTF

**アーカイブ**
- ZIP、RAR、7Z、TAR、GZ

**コードファイル**
- JS、TS、JSON、XML、HTML、CSS
- PY、GO、JAVA、CPP、C、SH、MD

### ファイル機能
- **アップロード方法**: ボタンクリック、ドラッグ&ドロップ、複数ファイル同時アップロード
- **ファイル管理**: アップロードファイルの一覧表示、ダウンロード、削除
- **ページ関連付け**: ファイルを特定のページに関連付けて管理
- **フィルタリング**: ファイルタイプ別の絞り込み表示
- **セキュリティ**: MIMEタイプ検証、ファイル名サニタイズ

### 制限事項
- 最大ファイルサイズ: 50MB
- 対応ファイル形式: 上記のリストに記載されたもの

## 🤝 貢献

プルリクエストやIssueは大歓迎です！

---

**Built with ❤️ using Go, Next.js, and real-time collaboration technologies.**
