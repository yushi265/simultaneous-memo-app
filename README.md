# simultaneous-memo-app

Notionライクなリアルタイム同時編集メモアプリケーション

## 技術スタック

### フロントエンド
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- TipTap (リッチテキストエディター)
- Yjs (リアルタイム同期)

### バックエンド
- Go + Echo v4
- PostgreSQL
- WebSocket (リアルタイム通信)

## 開発環境のセットアップ

```bash
# Dockerで起動
docker-compose up

# フロントエンド開発サーバー (http://localhost:3000)
cd frontend && npm run dev

# バックエンド開発サーバー (http://localhost:8080)
cd backend && air
```
