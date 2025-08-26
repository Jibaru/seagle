# Seagle

AI-powered PostgreSQL database management tool built with Wails (Go + React/TypeScript).

## Project Overview

Seagle is a desktop application similar to JetBrains DataGrip, designed to connect to PostgreSQL databases and perform queries using AI assistance.

## 📋 Features Status

### ✅ Database Connection - **COMPLETED**
- **Connection Methods:**
  - Form-based connection (host, port, database, username, password, SSL mode)
  - Connection string support for advanced configurations
  - SSL mode options (disable, require, verify-ca, verify-full)
- **Connection Management:**
  - Test connection functionality before connecting
  - Connect/disconnect with proper state management
  - Connection status tracking
  - Saved connections with persistent storage (JSON file)
  - Connect by ID for quick access to saved connections
- **Database Discovery:**
  - Automatic database enumeration after successful connection
  - PostgreSQL system database filtering
  - Alphabetically sorted database list

### ✅ Database Schema Explorer - **COMPLETED**
- **Sidebar Navigation:**
  - Hierarchical tree view (Databases → Tables → Columns)
  - Expandable/collapsible database and table nodes
  - Visual loading states for async operations
  - Database, table, and column selection with highlighting
- **Table Structure:**
  - Automatic table listing for selected databases
  - Column information display (name, data type, nullable, default values)
  - Real-time loading indicators

### ✅ Query Interface - **COMPLETED**
- **SQL Editor:**
  - Full-featured SQL query editor
  - Query execution with keyboard shortcuts
  - Database context awareness
  - Execute/Stop query functionality
- **Results Display:**
  - Tabular results presentation
  - Query performance metrics (duration, rows affected)
  - Error handling and display
  - Support for both SELECT and DML/DDL operations

### ✅ User Interface - **COMPLETED**
- **Modern Layout:**
  - DataGrip-inspired interface design
  - Dark/Light theme toggle support
  - Responsive sidebar and main content areas
  - Header with database/table context display
- **State Management:**
  - React Context for theme management
  - Zustand stores for database and connection state
  - Persistent UI state across sessions

## 🏗️ Current Architecture

### Directory Structure
```
seagle/
├── main.go                    # Application entry point with handler bindings
├── app.go                     # Minimal Wails app struct (context only)
├── go.mod                     # Go module dependencies
├── wails.json                 # Wails project configuration
├── Makefile                   # Build automation
├── README.md                  # Standard Wails template README
├── CLAUDE.md                  # Project documentation and architecture
├── core/                      # Core business logic (Clean Architecture)
│   ├── domain/                # Domain entities and repository interfaces
│   │   ├── connection.go      # Connection domain entity with business rules
│   │   └── connection_repo.go # Repository interface
│   ├── infra/                 # Infrastructure layer
│   │   ├── handlers/          # Wails handlers (application layer)
│   │   │   ├── connect.go         # Database connection handler
│   │   │   ├── connect_by_id.go   # Connect by saved ID handler
│   │   │   ├── disconnect.go      # Disconnection handler
│   │   │   ├── test_connection.go # Connection testing handler
│   │   │   ├── get_tables.go      # Tables listing handler
│   │   │   ├── get_table_columns.go # Table columns handler
│   │   │   ├── execute_query.go   # Query execution handler
│   │   │   └── list_connections.go # Saved connections handler
│   │   └── persistence/       # Data persistence layer
│   │       ├── common.go          # Common utilities
│   │       └── connection_repo.go # JSON file repository implementation
│   └── services/              # Application services (use cases)
│       ├── connection.go      # Connection service with business logic
│       └── types/             # Shared type definitions
│           └── connection.go  # DTOs and data structures
└── frontend/                  # React/TypeScript frontend
    ├── package.json           # Frontend dependencies and scripts
    ├── vite.config.ts         # Vite build configuration
    ├── tailwind.config.cjs    # Tailwind CSS configuration
    ├── postcss.config.cjs     # PostCSS configuration
    ├── biome.json             # Biome linter configuration
    ├── src/
    │   ├── components/        # React components
    │   │   ├── ui/            # shadcn/ui components
    │   │   │   ├── button.tsx # Button component
    │   │   │   ├── input.tsx  # Input component
    │   │   │   └── label.tsx  # Label component
    │   │   ├── DatabaseConnectionForm.tsx # Connection form component
    │   │   ├── WelcomeScreen.tsx          # Landing screen with saved connections
    │   │   ├── MainLayout.tsx             # Main application layout
    │   │   ├── Sidebar.tsx                # Database tree navigation
    │   │   ├── QueryInterface.tsx         # SQL query interface
    │   │   ├── SqlEditor.tsx              # SQL editor component
    │   │   ├── QueryResults.tsx           # Query results display
    │   │   ├── SavedConnections.tsx       # Saved connections management
    │   │   └── ThemeToggle.tsx            # Dark/Light theme toggle
    │   ├── contexts/          # React contexts
    │   │   └── ThemeContext.tsx           # Theme management context
    │   ├── store/             # State management
    │   │   ├── DatabaseStore.tsx          # Database and UI state (Zustand)
    │   │   └── ConnectionsStore.tsx       # Connections state (Zustand)
    │   ├── lib/
    │   │   └── utils.ts       # Utility functions
    │   ├── App.tsx            # Main application component
    │   ├── App.css            # Global styles
    │   └── main.tsx           # React entry point
    └── wailsjs/               # Generated Wails bindings
        ├── go/
        │   ├── handlers/      # TypeScript bindings for Go handlers
        │   └── models.ts      # Generated type definitions
        └── runtime/           # Wails runtime bindings
```

## 🛠️ Technology Stack

### Backend (Go)
- **Framework**: Wails v2.10.2 (Desktop application framework)
- **Language**: Go 1.23
- **Database Driver**: github.com/lib/pq v1.10.9 (PostgreSQL)
- **Architecture Pattern**: Clean Architecture (Domain/Infrastructure/Services)
- **UUID Generation**: github.com/google/uuid v1.6.0 (Connection IDs)
- **Persistence**: JSON file-based storage for connection settings

### Frontend (React/TypeScript)
- **Framework**: React 18.2.0 with TypeScript 4.6.4
- **Build Tool**: Vite 3.0.7
- **UI Library**: shadcn/ui components based on Radix UI
- **Styling**: Tailwind CSS 3.4.17 with tailwindcss-animate
- **Icons**: Lucide React 0.541.0 (Bird icon for branding)
- **Linting**: Biome 1.9.4

### Key Dependencies

- **Backend (go.mod)**
- **Frontend (package.json)**

## 🎨 UI Design System

### Layout Structure (DataGrip-inspired)
```
┌─────────────────────────────────────────────────────┐
│ Welcome Screen with Seagle Logo (Bird Icon)        │
├─────────────────────────────────────────────────────┤
│ ┌─────────┐ ┌───────────────────────────────────────┐ │
│ │ sidebar │ │         Query Editor                  │ │
│ │ sidebar │ │                                       │ │  
│ │ sidebar │ │                                       │ │
│ │ sidebar │ │                                       │ │
│ │ sidebar │ ├───────────────────────────────────────┤ │
│ │ sidebar │ │     Results Table Format              │ │
│ └─────────┘ └───────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

### Color Scheme
- **Primary Background:** `rgba(27, 38, 54, 1)` (Dark blue-gray)
- **Component Background:** White cards with shadow
- **Primary Color:** Blue-600 (`bg-blue-600`)
- **Secondary Color:** Gray-200 (`bg-gray-200`)
- **Error Color:** Red-600 (`bg-red-600`)
- **Success Color:** Green-100 border with green-700 text

### Typography
- **Font Family:** "Nunito" (custom loaded), with system font fallbacks
- **Headings:** Bold weights with proper text hierarchy
- **Body:** Regular weight, high contrast for readability

## 📐 Architecture Rules

### Clean Architecture Implementation

#### Domain Layer (`/core/domain`)
- **Entities**: Core business objects (Connection) with encapsulated business rules
- **Repository Interfaces**: Contracts for data persistence (ConnectionRepo)
- **Business Logic**: Domain entities contain validation and business rules
- **No Dependencies**: Domain layer has no external dependencies

#### Infrastructure Layer (`/core/infra`)
- **Handlers** (`/core/infra/handlers`): Wails application endpoints
  - Single Responsibility: Each handler handles one specific action
  - DTO Pattern: Structured input/output with proper validation
  - Error Handling: Consistent error responses for frontend consumption
  - Naming Convention: `{Action}Handler` with `{Action}Input`/`{Action}Output`
- **Persistence** (`/core/infra/persistence`): Repository implementations
  - File-based JSON storage for connection configurations
  - Interface compliance with domain repository contracts

#### Service Layer (`/core/services`)
- **Use Cases**: Application-specific business logic orchestration
- **State Management**: Connection lifecycle and database operations
- **Cross-cutting Concerns**: Database connectivity, transaction management
- **Type Definitions** (`/core/services/types`): DTOs for data transfer

## 🔧 Build & Development

### Development Commands
- **Start Development**: `wails dev` (runs backend + frontend with hot reload)
- **Build Production**: `wails build` (creates executable)
- **Frontend Development**: `npm run dev` (Vite dev server in frontend/)
- **Frontend Build**: `npm run build` (TypeScript compile + Vite build)
- **Frontend Linting**: `npm run biome` (Biome check, lint, and format)
  - `npm run lint` (Biome linter with auto-fix)
  - `npm run format` (Biome formatter)
  - `npm run check` (Biome combined check)

### Configuration Files
- **wails.json**: Wails project configuration and build settings
- **go.mod**: Go module definition and dependencies
- **frontend/package.json**: npm scripts and dependencies
- **frontend/vite.config.ts**: Vite bundler configuration
- **frontend/tailwind.config.js**: Tailwind CSS configuration
- **frontend/biome.json**: Biome linter and formatter settings

## 🔗 Data Flow

### Connection Process
1. **Frontend:** User fills connection form or selects saved connection
2. **Handler:** `ConnectHandler.Connect()` or `ConnectByIDHandler.ConnectByID()` receives input
3. **Service:** `ConnectionService.Connect()` creates/retrieves domain connection
4. **Persistence:** Connection saved to JSON file with UUID
5. **Domain:** Connection entity handles PostgreSQL connection logic
6. **Service:** `ConnectionService.GetDatabases()` queries available databases
7. **Handler:** Returns success status, message, and database list
8. **Frontend:** Updates UI state and displays database tree

### Query Execution Flow
1. **Frontend:** User enters SQL query in editor
2. **Handler:** `ExecuteQueryHandler.ExecuteQuery()` receives query and database
3. **Service:** `ConnectionService.ExecuteQuery()` connects to specific database
4. **Domain:** Connection handles query execution (SELECT/DML/DDL)
5. **Service:** Formats results with metadata (columns, rows, duration)
6. **Handler:** Returns structured query results or error
7. **Frontend:** Displays results in tabular format with performance metrics

### Schema Discovery Flow
1. **Frontend:** User expands database/table in sidebar
2. **Handler:** `GetTablesHandler` or `GetTableColumnsHandler` called
3. **Service:** Connection service queries PostgreSQL system tables
4. **Database:** Queries `information_schema` for metadata
5. **Service:** Formats table/column information
6. **Frontend:** Updates tree view with loading states and data

## 🎯 Component Architecture

### Screen Management
- **WelcomeScreen:** Entry point with branding and saved connections list
- **DatabaseConnectionForm:** Two-mode form (individual fields vs connection string)
- **MainLayout:** Primary interface with sidebar and query/table views

### State Management
```typescript
type ScreenState = 'welcome' | 'connection' | 'connected';

// Zustand stores
interface DatabaseState {
  databases: string[];
  selectedDatabase: string | null;
  selectedTable: string | null;
  expandedDatabases: Set<string>;
  expandedTables: Set<string>;
  databaseTables: Record<string, string[]>;
  tableColumns: Record<string, TableColumn[]>;
  loadingTables: Set<string>;
  loadingColumns: Set<string>;
}

interface ConnectionsState {
  connections: ConnectionSummary[];
  connectingId: string | null;
  isLoading: boolean;
}
```

### Component Features
- **Hierarchical Navigation**: Tree view with expand/collapse functionality
- **Async State Management**: Loading indicators for database operations
- **Theme Support**: Dark/light mode with React Context
- **Real-time Updates**: Zustand state synchronization across components
- **Error Handling**: Comprehensive error states and user feedback

## 🚦 Development Guidelines

### Code Standards
- **Language**: All code and comments in English
- **Incremental Development**: Implement one feature at a time
- **File Management**: Prefer editing existing files over creating new ones
- **Documentation**: Update CLAUDE.md with significant changes
- **Error Handling**: Comprehensive error handling with user-friendly messages

### Current Implementation Details
- **Database Connection**: Full PostgreSQL connectivity with SSL support and persistence
- **Schema Explorer**: Complete database/table/column tree navigation
- **Query Interface**: SQL editor with execution and results display
- **UI State Management**: Zustand stores with React Context for theme
- **Type Safety**: TypeScript throughout with proper Go-to-TS bindings
- **Component Architecture**: shadcn/ui components with consistent styling
- **Error Feedback**: Comprehensive error handling and loading states

### Potential Future Enhancements
- **AI Integration**: Query assistance and optimization suggestions
- **Advanced Query Features**: Query history, favorites, and templates
- **Data Export**: CSV, JSON, and other format exports
- **Connection Profiles**: Connection grouping and organization
- **Performance Monitoring**: Query execution analysis and optimization tips

## 🎯 Current Status Summary

### ✅ Completed Features
- **Database Connectivity**: Full PostgreSQL connection management with SSL support
- **Connection Persistence**: JSON-based storage with UUID identification
- **Schema Navigation**: Complete database/table/column tree explorer
- **Query Interface**: SQL editor with execution, results, and performance metrics
- **UI/UX**: Modern interface with dark/light theme support
- **State Management**: Comprehensive state handling with Zustand and React Context
- **Error Handling**: Robust error states and user feedback

### 📋 Technical Achievements
- **Clean Architecture**: Domain-driven design with proper layer separation
- **Type Safety**: Full TypeScript coverage with Go-to-TS bindings
- **Modern Tooling**: Vite build system, Tailwind CSS, Biome linting
- **Component Library**: shadcn/ui with Radix UI primitives
- **Performance**: Async operations with loading states
- **Code Quality**: Consistent patterns and comprehensive error handling

### 🚀 Project Status
Seagle has reached a **feature-complete MVP state** as a PostgreSQL database management tool. The core functionality equivalent to basic DataGrip features is fully implemented:

- ✅ Connection management with persistence
- ✅ Database schema exploration 
- ✅ SQL query execution with results
- ✅ Modern desktop interface
- ✅ Clean, maintainable architecture

The project provides a solid foundation for future enhancements while maintaining excellent code quality and user experience standards.
