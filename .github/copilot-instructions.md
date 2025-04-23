# Copilot Instructions

## Warning

- If you find any differences between these instructions and the actual implementation, update these instructions accordingly.

## Language

- Match user's language (Japanese/English/etc.)
- Default to English if unclear

## Commands

- Build: `npx nx build <project>`
- Test: `npx nx test <project>`
- Lint: `npx nx lint <project>`
- Format: `npx nx format:write`

## Workspace Structure

- `/apps/api` - Go API service (contributor image generation)
- `/apps/webapp` - Angular frontend application
- `/apps/worker` - TypeScript background worker
- `/firebase` - Firebase configuration
- `/tools` - Build and utility scripts

## Coding Rules

- Go: Follow standard Go idioms and error handling patterns
- TypeScript: Use strict typing, avoid `any` type
- Angular: Use standalone components and signals API
- Tests: Write unit tests for all new functionality
- Commits: Use conventional commit format
- Documentation: All code comments and inline documentation must be written in English, regardless of the language used in instructions

## Dependencies Management

- Node.js: Use pnpm for package management
- Go: Use go modules for dependency management
- Angular: Follow Angular versioning policy
- Update dependencies: Run `pnpm update` to update npm packages
- Go dependencies: Run `go get -u ./...` to update Go modules
- Security updates: Prioritize security-related dependency updates
