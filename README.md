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


<div style="display:flex;gap:25px">
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fironpark%2Fteatime?ref=badge_large">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fironpark%2Fteatime.svg?type=large" style="min-width:300px">
    </a>
    <div>
        <div>
        <h3>Why MPL 2.0</h3>
        Because I wanted it to be used as freely as possible while ensuring that improvements would be shared with the community. As long as the library itself isn't modified, there's no obligation to disclose source code, and unlike GPL-family licenses, it doesn't impose additional constraints even in Go's static linking environment, making it well-suited for this purpose.
        </div>
       <div>
       <h3>I Found GPL License in the Report - Is it OK?</h3>
       These are from dependencies that generate Go code from C headers without actual linking. Pure Go implementations that reference libc headers for compatibility don't typically create licensing obligations, but consult legal counsel for specific concerns. See <a href="https://gitlab.com/cznic/libc/-/issues/31#note_1644264616">related discussion</a>.
       </div>
    </div>
</div>