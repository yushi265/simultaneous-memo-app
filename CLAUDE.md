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

## 画像アップロード機能

リアルタイム協調メモアプリに完全な画像機能を実装しました。

### 実装済み機能

**バックエンド**
- 包括的な画像バリデーション（ファイルタイプ、サイズ制限）
- 画像リサイズ・最適化（Goライブラリ `imaging` 使用）
- 画像メタデータのデータベース保存（寸法、ファイルサイズ等）
- レスポンシブ画像配信エンドポイント（`/api/img/*`）
- サムネイル自動生成
- 画像削除機能（ファイルシステム + DB）
- ページと画像の関連付け管理
- 孤立した画像のクリーンアップ機能
- 年月別ディレクトリ構造（`/uploads/images/YYYY/MM/`）

**フロントエンド**
- TipTap画像拡張機能の統合
- エディターツールバーの画像アップロードボタン
- ドラッグ&ドロップでの画像アップロード
- クリップボードからの画像ペースト機能
- リサイズ可能な画像コンポーネント（ドラッグハンドル付き）
- アップロード進捗インジケーター
- 画像読み込み状態とエラーハンドリング
- レスポンシブ画像表示

**リアルタイム同期**
- Yjsによる画像ノードのリアルタイム同期
- 複数ユーザー間での画像操作の共有

### API エンドポイント

画像関連のAPIエンドポイント：
- `POST /api/upload` - 画像アップロード
- `GET /api/img/*` - レスポンシブ画像配信
- `GET /api/images` - 画像一覧取得
- `GET /api/images/:id` - 特定画像の取得
- `DELETE /api/images/:id` - 画像削除
- `POST /api/admin/cleanup-images` - 孤立画像のクリーンアップ
