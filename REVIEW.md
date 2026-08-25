Reviewed commit 407c4aa on feature/auth. Verdict: request changes.

### Findings

1. [P1] Registration accepts an unlimited request body

internal/features/authentication/transport/http/create_user.go:35 decodes r.Body without http.MaxBytesReader or another limit.

Because registration is public, a client can send a very large password/string. The JSON decoder allocates it before service validation rejects passwords over 72 bytes, creating a memory-exhaustion risk.

Limit the body before decoding, preferably in the shared request decoder or middleware.

2. [P1] The down migration fails after registering a user

migrations/000006_optional_user_profile.down.sql:5 restores NOT NULL:

ALTER COLUMN weight_grams SET NOT NULL,
ALTER COLUMN birth_date SET NOT NULL,
ALTER COLUMN height_cm SET NOT NULL;

But registration creates users with all three fields set to NULL. Therefore, after the first registration, migrate-down cannot apply.

This migration is not safely reversible without a backfill policy. Either provide valid values before restoring constraints or explicitly treat/document it as irreversible.

3. [P2] Valid email addresses are rejected

internal/features/authentication/service/create_user.go:24 uses:

^[a-z0-9][a-z0-9.]*[a-z0-9]@...

It requires at least two characters before @ and permits only letters, digits and dots. Consequently, valid addresses such as these are rejected:

a@example.com
john+training@example.com

The domain section is also too permissive when the service is called outside HTTP validation—for example, leading hyphens are accepted.

The existing database constraint has similar restrictions, so changing this requires a new migration rather than editing applied migration 000002.

4. [P2] Password length has two different meanings

The HTTP DTO uses validate:"min=8", which counts Unicode characters:

Password string `validate:"required,min=8"`

The service counts bytes:

passwordLength := len([]byte(password))

For example, a six-character Cyrillic password can exceed eight bytes: the service considers it long enough, while transport rejects it.

Define the rule explicitly:

- minimum in Unicode characters;
- bcrypt maximum in bytes (72).

Keeping the business rule in one place would prevent transport and service validation from diverging.

5. [P2] PostgreSQL error mapping discards the original error chain

internal/core/repository/postgres/pool/pgx/adapters.go:60 formats the original error with %v and wraps only ErrUnknown:

fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)

Because this commit also starts applying mapErrors to Query and Exec, callers can no longer detect errors such as:

errors.Is(err, context.DeadlineExceeded)
errors.Is(err, context.Canceled)

Preserve both the adapter category and original cause, for example with a custom wrapper or errors.Join.

### What is correct

- Passwords are hashed with bcrypt and never passed to PostgreSQL as plaintext.
- Email normalization happens before persistence.
- Duplicate email violations are mapped to domain ErrConflict.
- Layer boundaries are respected.
- SQL is parameterized.
- BIGINT version maps to int64.
- Tests cover hashing, normalization, invalid inputs and conflict propagation.

Verification passed:

go test ./...  — passed
go vet ./...   — passed

The largest missing test areas are HTTP handler behavior, PostgreSQL integration, request-size limits, email boundary cases, Unicode passwords and migration rollback. No code was changed.