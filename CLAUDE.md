# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

リアルタイムメモ - Real-time collaborative memo application with simultaneous editing capabilities and multi-user workspace support.

## Technology Stack

### Frontend
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS v3 for styling
- TipTap v2 for rich text editing with rich formatting options
- Yjs for real-time synchronization
- Zustand for state management with persistence
- Radix UI icons
- Custom logo component

### Backend
- Go 1.23 with Echo v4 framework
- PostgreSQL 16 with JSONB for page content and UUID for primary keys
- JWT authentication with bcrypt password hashing
- WebSocket for real-time communication with Gorilla WebSocket
- GORM v2 as ORM with datatypes and hooks support
- Rate limiting middleware for API protection
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
│   ├── components/    # React components (Header, Sidebar, Editor, Logo, FileUpload)
│   ├── lib/          # Utility functions (store, API client)
│   └── public/       # Static assets (logo.svg)
├── backend/           # Go API server
│   ├── config/       # Configuration management
│   ├── models/       # Database models (Page, BlockContent, Image, File)
│   ├── handlers/     # HTTP request handlers (pages, files, images)
│   └── websocket/    # WebSocket handlers (hub, client)
├── uploads/          # File upload directory
│   ├── images/      # Image files (YYYY/MM structure)
│   └── files/       # General files (YYYY/MM structure)
└── docker-compose.yml # Docker development environment
```

## Key Features

1. **User Authentication & Authorization**: JWT-based secure login with bcrypt password hashing
2. **Multi-Workspace Support**: Personal and team workspaces with role-based access control
3. **Real-time Collaborative Editing**: Using Yjs CRDT for conflict-free synchronization
4. **Rich Text Editor**: TipTap with headings, lists, code blocks, and rich formatting
5. **WebSocket Synchronization**: Real-time updates with user cursors and authentication
6. **PostgreSQL with JSONB**: Flexible content storage with UUID primary keys
7. **Performance Optimization**: Request caching, retry logic, and rate limiting
8. **Image Management**: Upload, resize, optimization, and responsive serving
9. **File Management**: General file upload with type validation and metadata storage
10. **Auto-save**: 3-second debounced saving with error handling
11. **Japanese UI**: Complete Japanese localization
12. **Docker Environment**: Containerized development setup

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration with auto workspace creation
- `POST /api/auth/login` - User login with JWT token
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user info and workspaces

### Workspaces
- `GET /api/workspaces` - List user workspaces
- `POST /api/workspaces` - Create new workspace
- `GET /api/workspaces/:id` - Get workspace details
- `PUT /api/workspaces/:id` - Update workspace
- `DELETE /api/workspaces/:id` - Delete workspace
- `POST /api/workspaces/:id/switch` - Switch to workspace (new JWT)

### Pages
- `GET /api/pages` - List workspace pages
- `POST /api/pages` - Create new page in current workspace
- `GET /api/pages/:id` - Get specific page
- `PUT /api/pages/:id` - Update page
- `DELETE /api/pages/:id` - Delete page

### Images
- `POST /api/upload` - Upload image with page association
- `GET /api/img/*` - Responsive image serving with optimization
- `GET /api/images` - List images in current workspace
- `GET /api/images/:id` - Get specific image metadata
- `DELETE /api/images/:id` - Delete image
- `POST /api/admin/cleanup-images` - Cleanup orphaned images

### Files
- `POST /api/upload/file` - Upload general file with validation
- `GET /api/files` - List files with filtering support
- `GET /api/files/:id` - Get file metadata
- `DELETE /api/files/:id` - Delete file
- `GET /api/file/*` - Serve uploaded file with access control

### WebSocket
- `GET /ws/:pageId` - Real-time sync with authentication support

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

## ファイルアップロード機能

リアルタイム協調メモアプリに汎用的なファイルアップロード機能を実装しました。

### 実装済み機能

**バックエンド**
- 汎用ファイルアップロードハンドラー（file_general.go）
- ファイルタイプバリデーション（PDF、ドキュメント、アーカイブ、コードファイル）
- ファイルサイズ制限（50MB）
- ファイルメタデータのデータベース保存（Fileモデル）
- ファイル配信エンドポイント
- ファイル削除機能（ファイルシステム + DB）
- ページとファイルの関連付け管理
- 年月別ディレクトリ構造（`/uploads/files/YYYY/MM/`）

**フロントエンド**
- FileUploadコンポーネント（ドラッグ&ドロップ対応）
- エディターツールバーのファイルアップロードボタン
- ファイル一覧表示（アイコン付き）
- ファイルダウンロード・削除機能
- アップロード進捗表示
- エラーハンドリング

### サポートファイルタイプ
- ドキュメント: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, TXT, CSV, RTF
- アーカイブ: ZIP, RAR, 7Z, TAR, GZ
- コードファイル: JS, TS, JSON, XML, HTML, CSS, PY, GO, JAVA, CPP, C, SH, MD

## 追加予定の機能

### コラボレーション機能
- ユーザー認証・権限管理
- メンション機能（@ユーザー名）
- コメント・注釈機能
- 変更履歴・バージョン管理
- ページ共有（読み取り専用リンク）

### エディター拡張
- マークダウンショートカット
- 数式エディター（LaTeX）
- 図表・チャート作成
- コードシンタックスハイライト
- テンプレート機能

### 組織・検索
- フォルダ・タグ管理
- 全文検索
- フィルター・ソート機能
- お気に入り・ピン留め
- アーカイブ機能

### その他
- オフライン編集対応
- モバイルアプリ
- 外部サービス連携（Slack、Google Drive等）
- AI機能（要約、翻訳、文章校正）
- エクスポート機能（PDF、Word）

## 実装状況

### ✅ フェーズ1: ユーザー認証 (完了)

**バックエンド**
- [x] User/Workspace/WorkspaceMemberモデル（UUID対応）
- [x] ワークスペース分離によるデータベースマイグレーション
- [x] bcryptパスワードハッシュ化 + JWT認証
- [x] ワークスペースコンテキスト付き認証ミドルウェア
- [x] 認証API完全実装（register/login/logout/me）
- [x] 全エンドポイントへのワークスペース制約追加

**フロントエンド**
- [x] バリデーション付きログイン・登録ページ
- [x] localStorage永続化対応Zustand認証ストア
- [x] SSRハイドレーション対応AuthGuardコンポーネント
- [x] リトライ機能付きトークンベースAPIクライアント
- [x] 保護ルートと認証フロー

**パフォーマンス最適化**
- [x] リクエストキャッシュとデバウンス
- [x] 指数バックオフリトライ付きレート制限
- [x] 強制ログアウトなしの429エラーハンドリング
- [x] WebSocket認証統合

### ✅ フェーズ2: ワークスペース管理 (完了)

- [x] ワークスペースCRUD操作
- [x] 新JWT生成によるワークスペース切り替え
- [x] WorkspaceSwitcher UIコンポーネント
- [x] 新規ワークスペース作成用CreateWorkspaceModal
- [x] ロールベース権限付きワークスペース設定ページ
- [x] 個人・チームワークスペースの区別

### 🚧 フェーズ3: コラボレーション機能 (未実装)

**次の優先事項:**
- [ ] トークンベースメンバー招待システム
- [ ] ロールベース権限システム（owner/admin/member/viewer）
- [ ] メンバー管理UI
- [ ] 権限ベースUI制御
- [ ] リアルタイムコラボレーション権限チェック

### 🔮 フェーズ4: 高度な機能 (計画中)

- [ ] ページ共有とパブリックリンク
- [ ] バージョン履歴と復元機能
- [ ] ページテンプレートとテンプレートライブラリ
- [ ] 高度な検索・フィルタリング
- [ ] 外部連携（Slack、Google Drive）
- [ ] AI機能（要約、翻訳）
- [ ] エクスポート機能（PDF、Word、Markdown）
