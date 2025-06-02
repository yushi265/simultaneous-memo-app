# システムアーキテクチャ

リアルタイムメモアプリの全体アーキテクチャとデータフローを説明します。

## 全体アーキテクチャ図

```mermaid
graph TB
    subgraph "フロントエンド (Next.js 14)"
        UI[React Components]
        TipTap[TipTap Editor]
        Yjs[Yjs CRDT]
        Store[Zustand State]
        
        UI --> TipTap
        TipTap --> Yjs
        Yjs --> Store
    end
    
    subgraph "リアルタイム通信"
        WS[WebSocket Connection]
        YjsProvider[Yjs WebSocket Provider]
        
        Yjs <--> YjsProvider
        YjsProvider <--> WS
    end
    
    subgraph "バックエンド (Go + Echo)"
        API[REST API Handlers]
        WSHandler[WebSocket Handler]
        Models[GORM Models]
        
        API --> Models
        WSHandler --> Models
    end
    
    subgraph "データストレージ"
        DB[(PostgreSQL 16)]
        Files[File System<br/>/uploads]
        
        Models --> DB
        API --> Files
    end
    
    subgraph "開発環境"
        Docker[Docker Compose]
        Air[Air Hot Reload]
        NextDev[Next.js Dev Server]
        
        Docker --> Air
        Docker --> NextDev
    end
    
    %% Connections
    UI --> API
    WS --> WSHandler
    
    classDef frontend fill:#e1f5fe
    classDef backend fill:#f3e5f5
    classDef storage fill:#e8f5e8
    classDef realtime fill:#fff3e0
    
    class UI,TipTap,Yjs,Store frontend
    class API,WSHandler,Models backend
    class DB,Files storage
    class WS,YjsProvider realtime
```

## データフロー図

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant UI as React UI
    participant Editor as TipTap Editor
    participant Yjs as Yjs CRDT
    participant WS as WebSocket
    participant Backend as Go Backend
    participant DB as PostgreSQL
    
    %% 初期ページ読み込み
    User->>UI: ページアクセス
    UI->>Backend: GET /api/pages/:id
    Backend->>DB: SELECT * FROM pages
    DB-->>Backend: ページデータ
    Backend-->>UI: JSON レスポンス
    UI->>Editor: コンテンツセット
    
    %% リアルタイム接続
    UI->>Yjs: Yjs初期化
    Yjs->>WS: WebSocket接続 /ws/:pageId
    WS->>Backend: WebSocket接続
    
    %% リアルタイム編集
    User->>Editor: テキスト入力
    Editor->>Yjs: 変更イベント
    Yjs->>WS: 変更データ送信
    WS->>Backend: 変更データ中継
    Backend->>WS: 他のクライアントに配信
    WS->>Yjs: 変更データ受信
    Yjs->>Editor: コンテンツ更新
    
    %% 自動保存
    Note over Editor,DB: 1秒デバウンス後
    Editor->>Backend: PUT /api/pages/:id
    Backend->>DB: UPDATE pages SET content=?
    DB-->>Backend: 更新完了
    Backend-->>Editor: 保存完了
```

## コンポーネント構成

```mermaid
graph TD
    subgraph "Pages"
        HomePage[Home Page]
        PageDetail[Page Detail]
    end
    
    subgraph "Layout Components"
        Header[Header]
        Sidebar[Sidebar]
        Logo[Logo Component]
    end
    
    subgraph "Editor Components"
        Editor[Editor]
        MenuBar[Editor Menu Bar]
        Content[Editor Content]
    end
    
    subgraph "State Management"
        Store[Zustand Store]
        API[API Client]
    end
    
    subgraph "Real-time"
        YjsDoc[Yjs Document]
        WSProvider[WebSocket Provider]
        Cursors[Collaboration Cursors]
    end
    
    HomePage --> Header
    HomePage --> Sidebar
    PageDetail --> Header
    PageDetail --> Sidebar
    PageDetail --> Editor
    
    Header --> Logo
    Editor --> MenuBar
    Editor --> Content
    Editor --> Cursors
    
    Editor --> Store
    Store --> API
    Editor --> YjsDoc
    YjsDoc --> WSProvider
    
    classDef page fill:#e3f2fd
    classDef layout fill:#f1f8e9
    classDef editor fill:#fce4ec
    classDef state fill:#fff8e1
    classDef realtime fill:#f3e5f5
    
    class HomePage,PageDetail page
    class Header,Sidebar,Logo layout
    class Editor,MenuBar,Content editor
    class Store,API state
    class YjsDoc,WSProvider,Cursors realtime
```

## 技術スタック詳細

### フロントエンド
- **Next.js 14**: App Router、TypeScript、Tailwind CSS
- **TipTap v2**: リッチテキストエディター
- **Yjs**: CRDT（Conflict-free Replicated Data Type）
- **Zustand**: 軽量状態管理
- **Radix UI**: アイコンライブラリ

### バックエンド
- **Go 1.23**: 高性能バックエンド言語
- **Echo v4**: 軽量Webフレームワーク
- **GORM v2**: Go用ORM
- **Gorilla WebSocket**: WebSocket実装
- **Air**: ホットリロード開発ツール

### データベース
- **PostgreSQL 16**: メインデータベース
- **JSONB**: リッチコンテンツ保存用
- **File System**: ファイルアップロード保存先

### 開発・デプロイ
- **Docker Compose**: 開発環境構築
- **Git**: バージョン管理
- **GitHub Actions**: CI/CD（予定）

## セキュリティ考慮事項

1. **入力検証**: フロントエンド・バックエンド両方で実装
2. **CORS設定**: 適切なオリジン制限
3. **SQL インジェクション対策**: GORM の Safe Query 使用
4. **XSS対策**: TipTapのHTMLサニタイゼーション
5. **ファイルアップロード制限**: ファイルタイプ・サイズ制限

## パフォーマンス最適化

1. **リアルタイム同期**: Yjs CRDTによる効率的な差分同期
2. **自動保存**: デバウンス機能でAPIコール最適化
3. **JSONB活用**: PostgreSQLの高速JSON操作
4. **コンポーネント最適化**: React.memo、useMemo活用
5. **バンドル最適化**: Next.js の自動最適化機能

## 拡張性

このアーキテクチャは以下の拡張に対応できます：

1. **マルチテナント**: ユーザー・組織管理
2. **権限システム**: ページ単位のアクセス制御
3. **プラグインシステム**: TipTap拡張機能
4. **API拡張**: RESTful API の機能追加
5. **スケーリング**: マイクロサービス分割対応