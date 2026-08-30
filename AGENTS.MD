## Project

SoloLeveling is an educational Go HTTP API backed by PostgreSQL and pgx.

The project is used to learn:
- Go;
- PostgreSQL and SQL;
- HTTP API design;
- layered architecture;
- Git collaboration;
- testing and development practices.

Prefer explanations that help the user understand the implementation instead
of providing unexplained code.

## Architecture

- `internal/core` contains code shared across application features.
- `internal/features` contains application features.
- Each feature is divided into:
- `transport/http`;
- `service`;
- `repository/postgres`.
- Follow the existing `exercises` feature as the basic structural example.
- Dependencies flow in this direction:

    transport -> service -> repository

- Transport handles HTTP requests, request DTOs, response DTOs and status codes.
- Services implement use cases and business orchestration.
- Repositories handle persistence and map PostgreSQL models to domain entities.
- Domain types must not depend on HTTP or PostgreSQL packages.
- HTTP DTOs must not be passed into services or repositories.
- PostgreSQL models must not be returned from repositories as domain results.

## PostgreSQL

- Application features depend on the `core_postgres_pool.Pool` abstraction.
- Do not use concrete pgx types outside the PostgreSQL adapter unless required.
- Always use parameterized SQL arguments such as `$1`, `$2`, and `$3`.
- Never build SQL by concatenating user-provided values.
- Close query rows and check `rows.Err()`.
- Wrap repository errors with operation context.
- Convert `pgx.ErrNoRows` into an appropriate domain or core error.

## Database migrations

- Do not reset, delete, or recreate database data unless explicitly asked.
- Entity identifiers use PostgreSQL `UUID` and Go `uuid.UUID`.
- Mutable entities use a `BIGINT version` column for optimistic locking.
- Foreign-key types must match the referenced primary-key types.

## HTTP API

- Public API routes use the `/api/v1` prefix.
- Use plural resource names, for example `/workouts` and `/exercises`.
- Use the standard `net/http` package and the existing server abstractions.
- Parse and validate path, query and body values in the transport layer.
- Do not trust identifiers supplied by clients for authorization.
- Return errors through the existing HTTP response handler.
- Preserve existing API contracts unless the requested change updates them.

## Go style

- Follow existing package naming and file organization.
- Keep one use case per file when practical.
- Prefer clear names over abbreviations.
- Wrap errors using `%w` when callers need `errors.Is`.
- Run `gofmt` after modifying Go files.
- Run `go test ./...` after modifying Go code.
- Do not add production dependencies without explaining why they are needed.

## Git and working tree safety

- Preserve existing and unrelated user changes.
- Do not commit, push, merge, rebase, reset, or switch branches unless
explicitly asked.
- Before suggesting a Git operation, inspect the current branch and working
tree when possible.
- Never use destructive Git commands such as `git reset --hard` without
explicit approval.
- Do not edit another developer's work merely to avoid a merge conflict.
- When multiple developers work in parallel, prefer separate feature branches
and separate endpoint files.

## Educational interaction

- Explain unfamiliar architecture and language concepts in plain language.
- Distinguish required changes from optional improvements.
- When reviewing code, identify:
- what is correct;
- what is incorrect;
- why it matters;
- which change is recommended.
- Do not silently fix an issue during an explanation or review request.
- Only modify code when the user explicitly asks to implement or change it.
- If an approach is outdated or inconsistent, point it out and explain the
consequences.
- Prefer guiding questions or small examples when they help the user learn.
- Do not hide important correctness, security, or data-loss risks.

## Scope control

- Make the smallest change that satisfies the requested task.
- Avoid unrelated refactoring.
- Do not change API contracts, migrations, or public types unless required by
the task.
- Ask for clarification only when different reasonable choices would
materially change the implementation.

## Project commands

Start PostgreSQL:

    make env-up

Start PostgreSQL port forwarding:

    make env-port-forward

Apply migrations:

    make migrate-up

Run the API:

    make app-run

Run all Go tests:

    go test ./...

# 🎯 Response Verbosity & Communication Style

You are assisting an experienced developer. Adjust your verbosity dynamically based on the intent of the prompt.

## 1. Quick & Direct Mode (Default for specific questions)
**Trigger:** The user asks a specific technical question, makes a binary choice ("A or B?"), asks how to fix a specific error, or asks for a specific command.
* **Rule:** Answer IMMEDIATELY and DEFINITIVELY. 
* **Format:** 1-5 sentences maximum + code snippet if needed.
* **FORBIDDEN:** Do not explain basic concepts. Do not write "It depends" without immediately providing the most pragmatic default choice. Do not list pros/cons unless explicitly requested.

## 2. Deep Dive Mode (Default for broad/learning questions)
**Trigger:** The user asks to explain a concept ("How does JWT work?"), requests an architectural overview, asks for best practices, or uses words like "explain", "how to build", or "what is the flow".
* **Rule:** Provide a comprehensive, full-length explanation.
* **Format:** Use structured Markdown (headers, lists, tables). Break down the logic step-by-step, explain the "why", and provide complete code examples.

## 3. Explicit Overrides
* If the prompt starts with **"TLDR"** or **"Быстро"**: Force Quick Mode.
* If the prompt starts with **"Explain"** or **"Подробно"**: Force Deep Dive Mode.