# nodex-desktop

nodex-desktop is a remote desktop application built with [Wails](https://wails.io/), Go, and React. 
The project is structured as a monorepo using [Turborepo](https://turbo.build/repo).

## Project Structure

### Apps

- `apps/desktop`: The main Wails desktop application.
  - `frontend/`: The React frontend for the desktop app.
  - `engine/`: The Go backend engine handling remote desktop capabilities (capture, input, transport, session).

### Packages

- `packages/ui`: A shared React component library using [shadcn/ui](https://ui.shadcn.com/) and Tailwind CSS v4.
- `packages/eslint-config`: Shared ESLint configurations.
- `packages/typescript-config`: Shared TypeScript configurations.

## Setup & Development

This project uses **bun** as the package manager.

### Install Dependencies

```sh
bun install
```

### Developing the Desktop App

To run the Wails desktop app in development mode:

```sh
cd apps/desktop
wails dev
```

*(Note: The frontend scripts are configured to use `bun` within `wails.json`)*

### Shared UI Components

The `packages/ui` library provides pre-built, accessible React components generated via `shadcn/ui`. 
It has been upgraded and configured with **Tailwind CSS v4**.

When integrating the UI components into an application, ensure you wrap your application root with the necessary providers (e.g., `TooltipProvider`). See `docs/general/ui-setup.md` for more details.
