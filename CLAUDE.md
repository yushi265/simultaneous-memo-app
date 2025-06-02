# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

リアルタイムメモ - Real-time collaborative memo application with simultaneous editing capabilities.

## Technology Stack

### Frontend
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS v3 for styling
- TipTap v2 for rich text editing with rich formatting options
- Yjs for real-time synchronization
- Zustand for state management
- Radix UI icons
- Custom logo component

### Backend
- Go 1.23 with Echo v4 framework
- PostgreSQL 16 with JSONB for page content
- WebSocket for real-time communication with Gorilla WebSocket
- GORM v2 as ORM with datatypes support
- Air for hot reloading in development

## Development Commands

### Frontend
```bash
cd frontend
npm run dev      # Start development server (http://localhost:3000)
npm run build    # Build for production
npm run lint     # Run ESLint
```

### Backend
```bash
cd backend
air              # Start development server with hot reload (http://localhost:8080)
go test ./...    # Run tests
go build         # Build binary
```

### Docker
```bash
docker-compose up -d --build  # Start all services in background with rebuild
docker-compose down           # Stop all services
docker-compose logs frontend  # View frontend logs
docker-compose logs backend   # View backend logs
docker-compose restart frontend  # Restart specific service
```

## Project Structure

```
/
├── frontend/           # Next.js frontend application
│   ├── app/           # App Router pages and layouts
│   ├── components/    # React components (Header, Sidebar, Editor, Logo)
│   ├── lib/          # Utility functions (store, API client)
│   └── public/       # Static assets (logo.svg)
├── backend/           # Go API server
│   ├── config/       # Configuration management
│   ├── models/       # Database models (Page, BlockContent)
│   ├── handlers/     # HTTP request handlers (pages, files)
│   └── websocket/    # WebSocket handlers (hub, client)
├── uploads/          # File upload directory
└── docker-compose.yml # Docker development environment
```

## Key Features

1. Real-time collaborative editing using Yjs CRDT
2. Rich text editor with TipTap (headings, lists, code blocks, formatting)
3. WebSocket-based synchronization with user cursors
4. PostgreSQL with JSONB for flexible content storage
5. Docker-based development environment
6. Custom logo and Japanese UI
7. File upload functionality
8. Auto-save with debouncing (1-second delay)

## API Endpoints

### Pages
- `GET /api/pages` - List all pages
- `POST /api/pages` - Create new page
- `GET /api/pages/:id` - Get specific page
- `PUT /api/pages/:id` - Update page
- `DELETE /api/pages/:id` - Delete page

### Files
- `POST /api/upload` - Upload file
- `GET /api/files/:id` - Get uploaded file

### WebSocket
- `GET /ws/:pageId` - WebSocket endpoint for real-time sync

## URLs

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

## TODO: 画像アップロード機能

### バックエンドタスク

1. **画像アップロードAPI拡張**
   - [x] 画像バリデーション追加（ファイルタイプ、サイズ制限）
   - [x] 画像リサイズ・最適化の実装（Goライブラリ `imaging` 使用）
   - [x] 画像メタデータの保存（ファイル名、サイズ、寸法、アップロード日時）
   - [x] 適切なMIMEタイプでの画像配信エンドポイント作成
   - [x] 画像削除機能の追加

2. **データベーススキーマ更新**
   - [x] メタデータ付き `images` テーブルの作成
   - [x] ページとの外部キーリレーション追加
   - [x] 画像参照を含むページコンテンツ構造の更新

3. **ファイル管理**
   - [x] アップロードディレクトリ構造の整理（`/uploads/images/YYYY/MM/`）
   - [x] 孤立した画像のクリーンアップ実装
   - [x] 画像サムネイル生成機能の追加

### フロントエンドタスク

4. **TipTap画像拡張機能**
   - [ ] `@tiptap/extension-image` のインストールと設定
   - [ ] カスタム画像アップロードハンドラーの作成
   - [ ] ドラッグ&ドロップでの画像アップロード追加
   - [ ] クリップボードからの画像ペースト実装

5. **UIコンポーネント**
   - [ ] エディターツールバーに画像アップロードボタン作成
   - [ ] 画像アップロード進捗インジケーター追加
   - [ ] エディター内での画像リサイズハンドル実装
   - [ ] 画像ギャラリー・ブラウザモーダル作成
   - [ ] 画像の代替テキスト入力ダイアログ追加

6. **エディター統合**
   - [ ] 画像ノードを扱うためのエディター更新
   - [ ] 画像配置の実装（インライン、ブロック、フロート）
   - [ ] 画像キャプション機能の追加
   - [ ] 画像読み込み状態とエラーの処理

### リアルタイム同期

7. **Yjs統合**
   - [ ] クライアント間での画像ノード同期の確認
   - [ ] 共同編集時の画像アップロード競合処理
   - [ ] リアルタイムシナリオでの画像操作テスト

### 技術的改善

8. **パフォーマンス＆UX**
   - [ ] 画像遅延読み込みの追加
   - [ ] プログレッシブ画像読み込みの実装
   - [ ] クライアントサイドでの画像圧縮追加
   - [ ] レスポンシブな画像表示の作成

9. **エラーハンドリング**
   - [ ] アップロード失敗時の適切なエラーメッセージ追加
   - [ ] アップロード中のネットワーク中断処理
   - [ ] 失敗時のリトライ機構実装

### テスト・ドキュメント

10. **品質保証**
    - [ ] 画像アップロード機能の単体テスト追加
    - [ ] 各種画像フォーマットでのテスト（JPEG、PNG、WebP、GIF）
    - [ ] ファイルサイズ制限とバリデーションのテスト
    - [ ] 画像エンドポイントのAPIドキュメント更新

### セキュリティ・検証

11. **セキュリティ対策**
    - [ ] クライアントとサーバー両方でのファイルタイプ検証実装
    - [ ] アップロードファイルのウイルススキャン追加（オプション）
    - [ ] パストラバーサル防止のためのファイル名サニタイズ
    - [ ] 画像アップロードのレート制限追加

### Docker・デプロイ

12. **インフラストラクチャ**
    - [ ] 永続的な画像保存のためのDockerボリューム更新
    - [ ] 静的画像配信用のnginx設定（本番環境）
    - [ ] アップロード設定用の環境変数追加
