# GEMINI.md

このファイルは、このリポジトリのコードを扱う際にGeminiが参照するためのガイダンスを提供します。

## プロジェクト概要

これは、複数ユーザーによる同時編集機能とマルチワークスペースをサポートする、Notionライクなリアルタイム共同編集メモアプリケーションです。

## 技術スタック

### フロントエンド
- **フレームワーク**: Next.js 14 (App Router使用)
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS v3
- **エディター**: TipTap v2（リッチテキスト編集用）
- **リアルタイム同期**: Yjs（CRDTベースの同期）
- **状態管理**: Zustand（永続化対応）
- **アイコン**: Radix UI

### バックエンド
- **言語/フレームワーク**: Go 1.23 + Echo v4
- **データベース**: PostgreSQL 16 (ページコンテンツにJSONB、主キーにUUIDを使用)
- **ORM**: GORM v2
- **認証**: JWT (パスワードハッシュ化はbcrypt)
- **リアルタイム通信**: Gorilla WebSocket
- **画像処理**: `disintegration/imaging`
- **ホットリロード**: Air

### インフラストラクチャ
- **コンテナ化**: Docker と Docker Compose

## 開発コマンド

### Docker利用（推奨）
- **全サービス起動**: `docker-compose up -d --build`
- **全サービス停止**: `docker-compose down`
- **ログ表示**: `docker-compose logs frontend` または `docker-compose logs backend`

### 各サービスの個別実行
- **フロントエンド**: `cd frontend && npm install && npm run dev` (http://localhost:3000 で実行)
- **バックエンド**: `cd backend && go mod download && air` (http://localhost:8080 で実行)
- **PostgreSQL**: `docker run --name postgres -e POSTGRES_PASSWORD=dev123 -e POSTGRES_DB=notion_app -p 5432:5432 -d postgres:16`

## プロジェクト構造

```
/
├── frontend/           # Next.js フロントエンドアプリケーション
│   ├── app/            # App Router ページとレイアウト
│   ├── components/     # Reactコンポーネント (Header, Sidebar, Editorなど)
│   ├── lib/            # ユーティリティ関数 (APIクライアント, Zustandストア)
│   └── public/         # 静的アセット
├── backend/            # Go APIサーバー
│   ├── config/         # 設定管理
│   ├── handlers/       # HTTPリクエストハンドラ
│   ├── middleware/     # Echoミドルウェア (認証, レート制限)
│   ├── models/         # GORMデータベースモデル
│   └── websocket/      # WebSocketハンドラ
├── uploads/            # ユーザーがアップロードしたファイルのディレクトリ
│   ├── images/         # 画像ファイル (YYYY/MM 形式で構造化)
│   └── files/          # 一般ファイル (YYYY/MM 形式で構造化)
├── docs/               # プロジェクトドキュメント
└── docker-compose.yml  # Docker開発環境
```

## 主な機能

- **ユーザー認証**: JWTによる安全なユーザー登録とログイン
- **ワークスペース管理**: ユーザーは複数のワークスペースを作成、切り替え、管理可能
- **リアルタイム共同編集**: YjsとWebSocketにより、複数ユーザーが同じドキュメントを同時に編集可能
- **リッチテキストエディター**: 画像やファイルのアップロードを含む、さまざまな書式設定オプションをサポートする高機能エディター
- **画像・ファイル処理**:
    - ドラッグ＆ドロップ、ペースト、またはボタンクリックによるアップロード
    - サーバーサイドでの画像リサイズと最適化
    - アクセス制御付きの安全なファイル配信
- **自動保存**: 過度なAPI呼び出しを防ぐためのデバウンス付きで変更を自動保存

## APIエンドポイント

主要なAPIエンドポイントの概要です。

### 認証
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/me`

### ワークスペース
- `GET /api/workspaces`
- `POST /api/workspaces`
- `PUT /api/workspaces/:id`
- `POST /api/workspaces/:id/switch`

### ページ
- `GET /api/pages`
- `POST /api/pages`
- `PUT /api/pages/:id`

### ファイル/画像アップロード
- `POST /api/upload` (画像)
- `POST /api/upload/file` (一般ファイル)
- `GET /api/img/*` (リサイズされた画像配信)
- `GET /api/file/*` (一般ファイル配信)

### WebSocket
- `GET /ws/:pageId`

## 貢献方法

1.  **アーキテクチャの理解**: `docs/architecture.md` と `docs/database-schema.md` を確認してください。
2.  **チェックリストに従う**: プルリクエストを作成する前に `DEVELOPMENT_CHECKLIST.md` を使用してください。
3.  **ドキュメントの更新**: 機能追加やAPI変更の際には、`README.md` と `CLAUDE.md` を更新してください。`pre-commit`フックがこれを思い出させるのに役立ちます。
4.  **コード品質の維持**: フロントエンドで `npm run lint` を実行し、Goのコードが正しくフォーマットされていることを確認してください。
5.  **テストの作成**: 新しいバックエンド機能には単体テストと統合テストを追加してください。

## 今後の作業（優先事項）

`CLAUDE.md`によると、次の優先事項はコラボレーション機能です:
1.  **メンバー招待システム**: トークンベースの招待機能
2.  **ロールベースアクセス制御 (RBAC)**: `owner`, `admin`, `member` のようなロールの実装
3.  **メンバー管理UI**: ワークスペースのメンバーを管理するためのUI