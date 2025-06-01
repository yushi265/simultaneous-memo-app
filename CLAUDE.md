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
