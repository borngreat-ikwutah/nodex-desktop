<!-- BEGIN:turborepo-agent-rules -->

# This is NOT the Turborepo you know

Turborepo configuration, task behavior, and CLI commands can vary between installed versions and may differ from your training data. Resolve the `turbo` package from this file's directory or relevant workspace; in monorepos, it may not be visible from the repository root. For example, run `node -p "require.resolve('turbo/package.json')"` from a workspace that depends on `turbo`.

Read `docs/README.md` inside that installed package first, then read the relevant pages from its `docs/` directory before changing Turborepo configuration or commands. Heed deprecation notices. These bundled docs match the installed package version and are available without network access.

This block is written and re-added by `turbo` before repository-scoped commands when an AI agent is detected. In the Turborepo source repository, its template is defined in `crates/turborepo-cli/src/cli/agent_guidance.rs`. Removing the managed block while updates are enabled means a later qualifying invocation will add it again. Set `"agentGuidance": false` in the root `turbo.json` or `turbo.jsonc` to opt out; this does not remove an existing block. Keep the block committed with your work to avoid an uncommitted change on the next agent invocation.
<!-- END:turborepo-agent-rules -->

# Software Guardrails & Development Guidelines

This repository enforces strict code quality and performance guardrails. Agents and developers must adhere to the following rules when interacting with this codebase:

## 1. Package Management & Build System

- **Bun**: Use `bun` for all package management tasks (`bun install`, `bun add`). Do not use `npm`, `yarn`, or `pnpm`.
- **Turborepo**: Rely on Turborepo for running tasks across workspaces. Use `bun run build`, `bun run lint`, and `bun run check-types` from the root to leverage caching and parallel execution.
- **Caching**: Ensure that `turbo.json` outputs are correctly mapped. Vite applications output to `dist/**` and `build/**`, not `.next/**`.

## 2. Code Quality & Formatting

- **Linting**: The repository uses a combination of shared ESLint configs and `oxlint` for blazing-fast linting.
- **Formatting**: `prettier` is strictly enforced for all supported files (`.ts`, `.tsx`, `.md`, `.css`, etc.).
- **Pre-commit Hooks**: `husky` and `lint-staged` are active. Commits will automatically format staged files and run type-checks. **Do not bypass these hooks (`--no-verify`)** unless absolutely necessary.

## 3. TypeScript & Type Safety

- **Strict Types**: The codebase operates with strict TypeScript configurations.
- **`check-types` Script**: Every workspace (`apps/*`, `packages/*`) must include a `check-types` script (typically `"check-types": "tsc --noEmit"`) in its `package.json` so that Turborepo can validate it during CI.

## 4. CI/CD Pipeline

- **GitHub Actions**: The CI pipeline (`.github/workflows/ci.yml`) is pinned to specific Bun versions and uses `--frozen-lockfile` for deterministic dependency resolution.
- **Local Verification**: Before creating a PR, verify your changes locally by running:
  - `bun run format`
  - `bun run lint`
  - `bun run check-types`
  - `bun run build`

## 5. Frontend Architecture & Folder Structure

The frontend (`apps/desktop/frontend`) enforces a strict modular folder structure. Agents must respect and organize code accordingly to maintain separation of concerns:

- `src/components/`: Reusable, generic UI components (e.g., layouts, buttons) that are completely decoupled from domain logic.
- `src/features/`: Domain-specific modules (e.g., `features/landing`, `features/devices`). Each feature should encapsulate its own components, hooks, api, and utilities (e.g., `features/devices/components/`).
- `src/routes/`: TanStack Router file-based route definitions. These files should be minimal and primarily import components directly from `features/`.
- `src/data/`: Static data, mock information, or configuration objects (e.g., `navigation.ts`).
- `src/hooks/`: Global custom React hooks. Feature-specific hooks should reside within their respective `features/` directory.
- `src/lib/`: Global utilities, helpers, and configurations.
- `src/assets/`: Static assets like images, icons, and fonts.

## 6. Go Backend & Desktop Shell

- **Location**: all Go code lives in `apps/desktop`. The Wails entrypoint is `main.go`, the bound application surface is `app.go`, and engine packages live under `apps/desktop/engine/<engine>` (each engine is self contained: it owns its persistence, its domain types, and its Wails bound `Service`). `engine/{capture,domain,input,session,transport}` are reserved and still empty.

### Engine 01 — Identity

`apps/desktop/engine/identity` gives every installation a permanent cryptographic identity. It is split by concern: `manager.go` (disk lifecycle), `machine.go` (host binding), `id.go` (ID derivation), `service.go` (Wails facade).

- **First boot** (`manager.go` `generate`): 32 bytes from `crypto/rand` seed an ed25519 keypair, a v4 UUID is minted, and the Node ID is derived. The record is written to `os.UserConfigDir()/Nodex/identity/machine.key` as JSON at `0600` inside a `0700` directory. The file is the single source of truth: delete it and a new identity is minted.
- **Subsequent boot** (`manager.go` `restore`): the file is parsed, the binding is checked, the sealed seed is unsealed, and the ed25519 private key is rebuilt. The Node ID never changes across restarts, reboots, or network changes.
- **ID derivation** (`id.go`): `sha256(seed)`, first 4 bytes big endian, then `% 900_000_000 + 100_000_000` so the ID is always 9 digits and never zero. `FormatNodeID` renders `XXX XXX XXX`. Note the space holds only 900M values, so Engine 02 discovery must disambiguate peers by `uuid`/`public_key` rather than by Node ID alone.
- **Sealed at rest**: the private key is never written as plaintext. `manager.go` seals it with AES-256-GCM under a key derived via `crypto/hkdf` from the `machine.go` binding, salted per UUID and AAD-bound to it. Copying `machine.key` to another host yields `ErrBindingMismatch`, not a stolen identity.
- **Machine binding is input, never identity** (`machine.go`): a SHA-256 fingerprint over hostname, uid, and non-loopback MACs. Changing a NIC cannot re-mint a Node ID; it only makes the sealed seed unreadable, which surfaces as `ErrBindingMismatch` rather than silently rotating the device. Do not invert this relationship in later engines.
- **Frontend facade** (`service.go`): `Service` exposes `GetNodeID`, `GetIdentity`, `GetPublicKey`, `KeyPath`. It returns `identity.Public`, never `identity.Identity`, so the seed can never cross the Wails bridge. `main.go` calls `identity.NewService()` before `wails.Run` and injects it through `NewApp`.
- `SigningKey()` rebuilds the ed25519 private key for Engines 03/04 transport handshakes.
- **Tests**: `go test ./engine/...` from `apps/desktop` covers first boot generation, `0600` enforcement, seed recovery across boots, no plaintext seed in the file, binding mismatch rejection, tampered ciphertext rejection, ID determinism and range, and corrupt file handling.

### Conventions

- **Bindings**: never hand edit `frontend/wailsjs`. Run `wails generate module` from `apps/desktop` after changing the bound surface. Only bind structs that are safe to serialise into the webview; add an explicit public DTO when a type holds secrets.
- **Storage**: identity lives in the OS config dir, not a portable `Nodex Data/` folder beside the binary. Keep it per-user so an identity never travels with a copied install. No SQLite layer exists yet.

## 7. Frontend Console (Home)

`features/home/` is a self-contained dark OLED control surface, decoupled from the shadcn card system on purpose.

- `home-page.tsx` assembles an asymmetric bento on a 12 column grid; `components/reveal.tsx` handles staggered `IntersectionObserver` entry on transform/opacity only, so nothing animates a layout property.
- `components/console-shell.tsx` is the mandatory double-bezel wrapper (outer tray `rounded-[2rem]` + recessed `rounded-[calc(2rem-0.375rem)]` core). Every major card nests through it. Do not add flat bordered panels.
- `components/mesh-backdrop.tsx` owns the fixed orbs and grain veil. It is fixed and `pointer-events-none` so scrolling content never repaints a blur.
- `use-node-identity.ts` is the only bridge to Go. It calls the generated bindings once on mount and guards against state updates after unmount.
- Motion uses `cubic-bezier(0.32, 0.72, 0, 1)` at 500-900ms. Never `linear` or `ease-in-out`. Interactive icons are ultra-light stroke (Tabler at `stroke={1.4}`).
