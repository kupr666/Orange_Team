### 5. Medium: login is vulnerable to brute force and timing-based user discovery

internal/features/authentication/service/login.go:41 returns immediately when the email does not exist, but performs an expensive bcrypt comparison for an existing email.

The public response is correctly normalized to “unauthorized,” but response timing can still reveal whether an account exists. There is also no rate limiting.

Recommended approach:
- Compare against a fixed dummy bcrypt hash when the user is absent.
- Add login rate limiting by IP and, optionally, normalized account identifier.
- Keep the same public response for unknown email and incorrect password.

### 6. Medium: browser clients will not get CORS support

The application reads HTTP_ALLOWED_ORIGINS, and a CORS middleware exists, but cmd/api/
main.go:83 never installs it.

A separate web frontend will have its login/register JSON requests blocked during the
browser’s preflight request. Add:
core_http_middleware.CORS(httpConfig.AllowedOrigins), to the global middleware chain.

### 7. Medium: string-based roles are duplicated across layers
The "user" and "admin" values appear independently in:

- the PostgreSQL constraint;
- JWT claims;
- Principal.Role;
- internal/features/exercises/transport/http/transport.go:51.