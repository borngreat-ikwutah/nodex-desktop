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
