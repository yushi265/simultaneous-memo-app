# リアルタイムメモ

複数ユーザーが同時にリアルタイムで編集可能なNotionライクなメモアプリケーション

![Logo](frontend/public/logo.svg)

## ✨ 主な機能

- 📝 **リッチテキストエディター**: 見出し、リスト、コードブロック、太字・斜体などの豊富なフォーマット
- 🔄 **リアルタイム同期**: 複数ユーザーが同時編集可能（Yjs CRDT使用）
- 👥 **協調カーソル**: 他のユーザーのカーソル位置をリアルタイム表示
- 💾 **自動保存**: 1秒のデバウンスで自動保存
- 🖼️ **画像機能**: ドラッグ&ドロップ、クリップボード、ボタンからの画像アップロード
- 🎛️ **画像編集**: エディター内でのリサイズ、レスポンシブ配信、サムネイル自動生成
- 🌐 **日本語対応**: 完全日本語化されたUI

## 🛠 技術スタック

### フロントエンド
- **Next.js 14** (App Router)
- **TypeScript** - 型安全性
- **Tailwind CSS v3** - スタイリング
- **TipTap v2** - リッチテキストエディター（画像拡張含む）
- **Yjs** - リアルタイム同期（CRDT）
- **Zustand** - 状態管理
- **Radix UI** - アイコン

### バックエンド
- **Go 1.23** + **Echo v4** - API フレームワーク
- **PostgreSQL 16** - データベース（JSONB使用）
- **GORM v2** - ORM
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

1. ブラウザで <http://localhost:3000> にアクセス
2. 「新規ページ」ボタンをクリックしてページを作成
3. タイトルと本文を編集
4. **画像アップロード**:
   - ツールバーの画像ボタンをクリック
   - エディターに画像をドラッグ&ドロップ
   - クリップボードから画像をペースト（Ctrl+V/Cmd+V）
5. **画像編集**: 画像を選択してリサイズハンドルでサイズ調整
6. 複数のブラウザタブを開いて同時編集をテスト

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

### 画像・ファイル管理
- `POST /api/upload` - 画像アップロード（ページID関連付け対応）
- `GET /api/img/*` - レスポンシブ画像配信（サムネイル対応）
- `GET /api/images` - 画像一覧取得
- `GET /api/images/:id` - 特定画像の詳細取得
- `DELETE /api/images/:id` - 画像削除
- `POST /api/admin/cleanup-images` - 孤立画像のクリーンアップ
- `GET /api/files/*` - 汎用ファイル取得

### リアルタイム通信
- `WebSocket /ws/:pageId` - リアルタイム同期

## 🏗 プロジェクト構造

```
├── frontend/                 # Next.js フロントエンド
│   ├── app/                 # App Router
│   ├── components/          # Reactコンポーネント
│   │   ├── Header.tsx       # ヘッダー
│   │   ├── Sidebar.tsx      # サイドバー
│   │   ├── Editor.tsx       # エディター（画像対応）
│   │   ├── EditorMenuBar.tsx # エディターツールバー
│   │   ├── ResizableImage.tsx # リサイズ可能画像コンポーネント
│   │   └── Logo.tsx         # ロゴ
│   ├── lib/                 # ユーティリティ
│   │   ├── image-upload.ts  # 画像アップロード処理
│   │   ├── image-utils.ts   # 画像関連ユーティリティ
│   │   └── image-resize-extension.ts # TipTap画像拡張
│   └── public/              # 静的ファイル
├── backend/                 # Go バックエンド
│   ├── config/              # 設定管理
│   ├── models/              # データモデル（画像テーブル含む）
│   ├── handlers/            # HTTPハンドラー
│   │   ├── image*.go        # 画像処理関連ハンドラー
│   │   └── file.go          # ファイルアップロード
│   └── websocket/           # WebSocket処理
├── uploads/                 # アップロードファイル
│   └── images/              # 画像ファイル（YYYY/MM構造）
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

## 🤝 貢献

プルリクエストやIssueは大歓迎です！

---

**Built with ❤️ using Go, Next.js, and real-time collaboration technologies.**
