# リアルタイムメモ

複数ユーザーが同時にリアルタイムで編集可能なNotionライクなメモアプリケーション

![Logo](frontend/public/logo.svg)

## ✨ 主な機能

- 📝 **リッチテキストエディター**: 見出し、リスト、コードブロック、太字・斜体などの豊富なフォーマット
- 🔄 **リアルタイム同期**: 複数ユーザーが同時編集可能（Yjs CRDT使用）
- 👥 **協調カーソル**: 他のユーザーのカーソル位置をリアルタイム表示
- 💾 **自動保存**: 1秒のデバウンスで自動保存
- 📎 **ファイルアップロード**: 画像やファイルのアップロード機能
- 🌐 **日本語対応**: 完全日本語化されたUI

## 🛠 技術スタック

### フロントエンド
- **Next.js 14** (App Router)
- **TypeScript** - 型安全性
- **Tailwind CSS v3** - スタイリング
- **TipTap v2** - リッチテキストエディター
- **Yjs** - リアルタイム同期（CRDT）
- **Zustand** - 状態管理
- **Radix UI** - アイコン

### バックエンド
- **Go 1.23** + **Echo v4** - API フレームワーク
- **PostgreSQL 16** - データベース（JSONB使用）
- **GORM v2** - ORM
- **Gorilla WebSocket** - リアルタイム通信
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

1. ブラウザで <http://localhost:3000> にアクセス
2. 「新規ページ」ボタンをクリックしてページを作成
3. タイトルと本文を編集
4. 複数のブラウザタブを開いて同時編集をテスト

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
```

## 📡 API エンドポイント

### ページ管理
- `GET /api/pages` - ページ一覧取得
- `POST /api/pages` - ページ作成
- `GET /api/pages/:id` - ページ詳細取得
- `PUT /api/pages/:id` - ページ更新
- `DELETE /api/pages/:id` - ページ削除

### ファイル管理
- `POST /api/upload` - ファイルアップロード
- `GET /api/files/:id` - ファイル取得

### リアルタイム通信
- `WebSocket /ws/:pageId` - リアルタイム同期

## 🏗 プロジェクト構造

```
├── frontend/                 # Next.js フロントエンド
│   ├── app/                 # App Router
│   ├── components/          # Reactコンポーネント
│   │   ├── Header.tsx       # ヘッダー
│   │   ├── Sidebar.tsx      # サイドバー
│   │   ├── Editor.tsx       # エディター
│   │   └── Logo.tsx         # ロゴ
│   ├── lib/                 # ユーティリティ
│   └── public/              # 静的ファイル
├── backend/                 # Go バックエンド
│   ├── config/              # 設定管理
│   ├── models/              # データモデル
│   ├── handlers/            # HTTPハンドラー
│   └── websocket/           # WebSocket処理
├── uploads/                 # アップロードファイル
└── docker-compose.yml       # Docker設定
```

## 📄 ライセンス

MIT License

## 🤝 貢献

プルリクエストやIssueは大歓迎です！

---

**Built with ❤️ using Go, Next.js, and real-time collaboration technologies.**
