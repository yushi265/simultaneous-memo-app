# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Notion-like real-time collaborative memo application with simultaneous editing capabilities.

## Technology Stack

### Frontend
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS for styling
- TipTap for rich text editing
- Yjs for real-time synchronization
- Zustand for state management

### Backend
- Go with Echo v4 framework
- PostgreSQL with JSONB for page content
- WebSocket for real-time communication
- GORM v2 as ORM

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
docker-compose up        # Start all services
docker-compose down      # Stop all services
docker-compose up -d     # Start in background
```

## Project Structure

```
/
├── frontend/           # Next.js frontend application
│   ├── app/           # App Router pages and layouts
│   ├── components/    # React components
│   └── lib/          # Utility functions and hooks
├── backend/           # Go API server
│   ├── models/       # Database models
│   ├── handlers/     # HTTP request handlers
│   └── websocket/    # WebSocket handlers
└── uploads/          # File upload directory
```

## Key Features

1. Real-time collaborative editing using Yjs CRDT
2. Rich text editor with TipTap
3. WebSocket-based synchronization
4. PostgreSQL with JSONB for flexible content storage
5. Docker-based development environment