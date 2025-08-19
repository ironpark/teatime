# Teatime
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fironpark%2Fteatime.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fironpark%2Fteatime?ref=badge_shield)


An AI-friendly local workflow automation program inspired by ComfyUI and n8n. Built with Wails 3, SvelteKit, and bits-ui for a modern desktop experience.
> [!WARNING]
> This project is in very early development stage. Core features are not yet properly implemented and may not work as expected.

## Build Requirements

- **Go** 1.24 or later
- **Node.js** 22.14.0 or later
- **pnpm** for package management (frontend)
## Getting Started
### Release
> Not ready yet.

### From Source


```bash
# Install Wails 3 (if not already installed):
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```


```bash
git clone https://github.com/ironpark/teatime && cd teatime
# Install Go dependencies
go mod tidy   
# Development mode
wails3 dev   
# Or build for production
wails3 build
```

## Project Structure

```
teatime/
├── frontend/
│   ├── src/
│   │   ├── routes/             # SvelteKit routes
│   │   ├── stories/            # Storybook stories
│   │   └── hooks.ts            # SvelteKit hooks
│   ├── bindings/               # Auto-generated Wails bindings
│   ├── messages/               # i18n message files
│   └── .storybook/             # Storybook configuration
│
├── services/                   # Go backend services
│
├── build/                      # Build configuration
│   ├── darwin/                 # macOS build config
│   ├── linux/                  # Linux build config
│   └── windows/                # Windows build config
│
├── main.go                     # Application entry point
└── go.mod                      # Go module definition
```

## Built With ❤

### Frontend:
- **[SvelteKit](https://svelte.dev/docs/kit/introduction)** - Web framework
- **[Bits-UI](https://bits-ui.com/)** - Headless components library for svelte
- **[shadcn-svelte](https://shadcn-svelte.com/)** - Modern UI components for Svelte
- **[Tailwind CSS v4](https://tailwindcss.com/)** - Utility-first CSS framework
- **[TypeScript](https://www.typescriptlang.org/)** - Type-safe JavaScript
- **[Storybook](https://storybook.js.org/)** - Component development environment

### Backend:
- **[Go](https://go.dev/)** - System programming language
- **[Wails 3](https://github.com/wailsapp/wails)** - Desktop application framework
- **[Extism Go SDK](https://github.com/extism/go-sdk)** - WebAssembly plugin system
- **[Sqlc](https://sqlc.dev/)** - Type-safe SQL code generator
- **[modernc.org/sqlite](https://modernc.org/sqlite)** - Pure Go SQLite implementation


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fironpark%2Fteatime.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fironpark%2Fteatime?ref=badge_large)