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
- **Database Discovery:**
  - Automatic database enumeration after successful connection
  - PostgreSQL system database filtering
  - Alphabetically sorted database list

### 🚧 Roadmap - **PENDING**
- List databases in the sidebar
- List tables of the database in the sidebar (like a tree view)
- List fields of each table in the sidebar  
- Perform queries on the selected database and show results in table format

## 🏗️ Current Architecture

### Directory Structure
```
seagle/
├── main.go                    # Application entry point
├── app.go                     # Minimal Wails app struct (context only)
├── go.mod                     # Go module dependencies
├── wails.json                 # Wails project configuration
├── gorules.md                 # Architecture rules and conventions
├── README.md                  # Standard Wails template README
├── core/                      # Core business logic
│   ├── handlers/              # Individual request handlers (Wails bindings)
│   │   ├── connect.go         # Database connection handler
│   │   ├── test_connection.go # Connection testing handler
│   │   └── disconnect.go      # Disconnection handler
│   ├── services/              # Business logic services
│   │   ├── connection.go      # Connection service with GetDatabases()
│   │   └── types/             # Shared type definitions
│   │       └── connection.go  # DatabaseConfig, DatabaseConnection types
└── frontend/                  # React/TypeScript frontend
    ├── package.json           # Frontend dependencies and scripts
    ├── vite.config.ts         # Vite build configuration
    ├── tailwind.config.js     # Tailwind CSS configuration
    ├── postcss.config.js      # PostCSS configuration
    ├── biome.json             # Biome linter configuration
    ├── src/
    │   ├── components/        # React components
    │   │   ├── ui/            # shadcn/ui components
    │   │   │   ├── button.tsx # Button component
    │   │   │   ├── input.tsx  # Input component
    │   │   │   └── label.tsx  # Label component
    │   │   ├── DatabaseConnectionForm.tsx # Main connection form
    │   │   └── WelcomeScreen.tsx          # Landing screen
    │   ├── App.tsx            # Main application component
    │   └── App.css            # Global styles
    └── wailsjs/               # Generated Wails bindings
        └── go/handlers/       # TypeScript bindings for Go handlers
```

## 🛠️ Technology Stack

### Backend (Go)
- **Framework**: Wails v2.10.2 (Desktop application framework)
- **Language**: Go 1.23
- **Database Driver**: github.com/lib/pq v1.10.9 (PostgreSQL)
- **Architecture Pattern**: Clean Architecture (Handlers/Services/Types)

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

## 📐 Architecture Rules (gorules.md)

### Handler Layer (`/core/handlers`)
- **Single Responsibility**: Each handler is a struct with exactly 1 method
- **DTO Pattern**: Each handler receives input as struct (DTO) if needed
- **Return Pattern**: Each handler returns output as struct pointer (DTO) and error if needed  
- **API Endpoints**: Handlers act as API endpoints for the frontend (Wails bindings)
- **Naming Convention**:
  - Handler: `{Action}Handler` (e.g., `ConnectHandler`)
  - Method: Same as action without "Handler" suffix (e.g., `Connect`)
  - Input: `{Action}Input` (e.g., `ConnectInput`)
  - Output: `{Action}Output` (e.g., `ConnectOutput`)

### Service Layer (`/core/services`)
- **Business Rules**: Services contain the core business logic of the application
- **State Management**: Can maintain state (desktop environment allows this)
- **Use Cases**: Single service can have multiple methods as use-case entry points
- **Data Access**: Responsible for database operations and external API calls

### Type Layer (`/core/services/types`)
- **Shared Structures**: Common data types used across handlers and services
- **Data Transfer**: Clean interfaces for data passing between layers
- **JSON Serialization**: Proper JSON tags for frontend communication

## 🔧 Build & Development

### Development Commands
- **Start Development**: `wails dev` (runs backend + frontend with hot reload)
- **Build Production**: `wails build` (creates executable)
- **Frontend Development**: `npm run dev` (Vite dev server in frontend/)
- **Frontend Build**: `npm run build` (TypeScript compile + Vite build)
- **Frontend Linting**: `npm run lint` (Biome linter with auto-fix)

### Configuration Files
- **wails.json**: Wails project configuration and build settings
- **go.mod**: Go module definition and dependencies
- **frontend/package.json**: npm scripts and dependencies
- **frontend/vite.config.ts**: Vite bundler configuration
- **frontend/tailwind.config.js**: Tailwind CSS configuration
- **frontend/biome.json**: Biome linter and formatter settings

## 🔗 Data Flow

### Connection Process
1. **Frontend:** User fills connection form or connection string
2. **Handler:** `ConnectHandler.Connect()` receives `ConnectInput` 
3. **Service:** `ConnectionService.Connect()` establishes PostgreSQL connection
4. **Service:** `ConnectionService.GetDatabases()` queries available databases
5. **Handler:** Returns `ConnectOutput` with success status, message, and database list
6. **Frontend:** Updates UI state and displays available databases

## 🎯 Component Architecture

### Screen Management
- **WelcomeScreen:** Entry point with Seagle branding and "New Connection" CTA
- **DatabaseConnectionForm:** Two-mode form (individual fields vs connection string)
- **Success Screen:** Post-connection confirmation with database info

### State Management
```typescript
type ScreenState = 'welcome' | 'connection' | 'connected';
```

### Form Features
- **Toggle Modes**: Radio buttons to switch between form fields and connection string
- **Real-time Validation**: Field-level validation with error display
- **Loading States**: Visual feedback during connection attempts
- **SSL Configuration**: Dropdown with PostgreSQL SSL modes

## 🚦 Development Guidelines

### Code Standards
- **Language**: All code and comments in English
- **Incremental Development**: Implement one feature at a time
- **File Management**: Prefer editing existing files over creating new ones
- **Documentation**: Update CLAUDE.md with significant changes
- **Error Handling**: Comprehensive error handling with user-friendly messages

### Current Implementation Details
- **Database Connection**: Full PostgreSQL connectivity with SSL support
- **UI State Management**: React state with proper screen transitions
- **Type Safety**: TypeScript throughout with proper Go-to-TS bindings
- **Component Architecture**: shadcn/ui components with consistent styling
- **Error Feedback**: User-friendly error messages and loading indicators

### Next Development Phase
- **Sidebar Implementation**: Database and table tree view
- **Schema Explorer**: Table and column information display
- **Query Interface**: SQL editor with syntax highlighting and execution
- **AI Integration**: Query assistance and optimization suggestions

## 🎯 Current Status Summary

### ✅ Completed Features
- PostgreSQL connection management (connect/test/disconnect)
- Dual-mode connection form (fields vs connection string)
- Database discovery and enumeration
- Clean architecture with proper separation of concerns
- Modern React/TypeScript frontend with shadcn/ui
- Responsive UI with loading states and error handling
- SSL connection support with multiple modes

### 📋 Technical Achievements
- Wails v2.10.2 desktop application framework integration
- Clean architecture following gorules.md specifications
- Individual handlers with proper DTOs and error handling
- PostgreSQL database service with connection pooling
- Modern frontend with Vite, Tailwind CSS, and TypeScript
- Component library with shadcn/ui and Radix UI primitives
- Biome linting for code quality and consistency

This architecture provides a solid foundation for building the remaining database exploration and query features while maintaining clean separation of concerns and established patterns.
