**Project API Flow — BriefBuletin_API

This document explains, in plain language, how the API starts and how requests flow through the codebase. It's written for a new developer (a fresher) so they can understand the project's structure, how authentication works, how a bypass or dev-mode works, and how to run the app locally.

**Project Layout**
- **`main.go`**: Entry point — starts configuration, database, router, and server.
- **`internal/service/`**: Business logic layer — each file implements a service (articles, comments, users, authentication, etc.).
- **`internal/dbmodel/`**: Database access layer — generated SQL query code and models (sqlc or similar). This is where low-level DB queries live.
- **`internal/model/`**: Application data models used across the services and handlers.
- **`util/`**: Reusable utilities: DB connection wrapper, AWS file access, common helpers.
- **`connection/dev-config.json`**, **`config.json`**: Configuration for environments (dev/test/prod). Use these for DB credentials, SMTP, AWS keys, and any feature flags like bypass.
- **`sqlc/`**: SQL schema and queries used by code generator (if present).

**High-level Flow (one sentence)**
- `main.go` loads config -> connects DB -> registers routes and middleware -> starts HTTP server. Incoming requests hit routes -> middleware (logging/auth) -> service layer -> DB layer -> response.

**Detailed Startup (what `main.go` typically does)**
- Load configuration: read `config.json` and environment-specific files (e.g., `connection/dev-config.json`).
- Initialize logger (if present).
- Open database connection using utilities in `util/dbconnectionwraper.go`.
- Initialize repositories / queriers from `internal/dbmodel/db_query` (these are the data-access functions).
- Wire services in `internal/service/` (pass DB or repo instances to services).
- Create HTTP router and register routes using `routeService.go` or similar.
- Attach global middleware (CORS, logging, authentication middleware).
- Start the HTTP server (listen on configured port). The compiled binary is in `build/exec/briefbuletin` after `make` or `go build`.

**Request Handling Flow (example: posting a comment)**
1. Client -> HTTP POST `/articles/{id}/comments` with JSON body.
2. Router matches route and calls the handler from `routeService.go` or `rest_Service.go`.
3. Middleware runs first: logging, request ID, and authentication middleware (if route requires auth).
4. Handler parses request into a model struct in `internal/model/comments.go` and performs basic validation.
5. Handler calls the `commentService` in `internal/service/commentService.go`.
6. `commentService` contains business rules (e.g., check article exists, sanitize content, enforce rate limits).
7. `commentService` calls the DB layer (`internal/dbmodel/db_query/querier.go`) to insert the comment.
8. DB returns result/ID -> service formats response -> handler returns HTTP response with JSON and proper status code.

**Service Layer Responsibilities**
- Implement business logic and rules.
- Accept input DTOs from handlers and map them to DB models.
- Call DB layer methods (querier) — do not hold database connection details in handlers.
- Return simple, typed results or domain errors (which handlers convert to HTTP errors).

**DB Layer (what `internal/dbmodel` does)**
- Contains generated query functions (by `sqlc` or handwritten queries).
- Exposes typed methods such as `CreateComment(ctx, params)` or `GetArticleByID(ctx, id)`.
- Services should call these methods rather than writing raw SQL.

**Authentication Flow (typical JWT-based approach)**
- Login flow:
  - Client POSTs credentials to `/auth/login`.
  - `authenticationService.go` validates credentials against user records in DB (`internal/model/userModel.go`).
  - On success, the service issues a JWT token (signed with a secret from `config.json`) and returns it to client.
- Protected routes:
  - A middleware extracts the token from the `Authorization: Bearer <token>` header.
  - Middleware validates token signature, expiry, and (optionally) checks claims (roles/user id).
  - If valid, the middleware attaches user information to the request context so handlers and services can read `ctx.Value("user")` or a typed context helper.
  - If invalid, middleware returns `401 Unauthorized`.

**Where authentication code usually lives**
- Token creation and validation functions: `internal/service/authenticationService.go` or a `util/jwt` helper.
- Handler endpoints: in `rest_Service.go` / `routeService.go`.
- Middleware: either a function in `routeService.go` or a dedicated `auth_middleware.go` that uses `authenticationService` to validate tokens.

**Bypass / Dev Mode (common patterns and how to find/use it)**
- Many projects include a dev-only bypass for local testing. Look for any of these indicators:
  - A `BYPASS_AUTH` or `DEV_MODE` flag in `config.json` or `connection/dev-config.json`.
  - An environment variable, e.g., `BYPASS_AUTH=true`.
  - Conditional code in `main.go` or in authentication middleware like `if cfg.DevMode { skip auth }`.
- If present, bypass mode commonly works like this:
  - When enabled, middleware does NOT validate JWTs and may inject a default test user into context.
  - This allows a developer to call protected endpoints without a token during development.
- Warning: never enable bypass in production. Always check production config or environment variables for any bypass flags before deploying.

**Configuration & Secrets**
- Keep secrets out of source control. Use environment variables or a secrets manager for DB credentials, JWT secret, and SMTP/AWS keys. The project has `connection/dev-config.json` for local dev config — do not use that for production credentials.
- Typical config values: DB host/port, DB user/password, JWT secret, JWT expiry, SMTP credentials, AWS S3 bucket/key, server port.

**How to Run Locally (quick start)**
- Ensure you have Go installed (Go 1.20+ recommended). If the project uses `Makefile`, use the Make targets.
- Example commands (PowerShell):
```pwsh
# build
go build -o build\exec\briefbuletin main.go
# run binary
./build/exec/briefbuletin
# or run directly with go
go run main.go
```
- If there's a `Makefile`, common targets are `make build` or `make run`.
- Ensure DB is running and `config.json` or `connection/dev-config.json` points to your local DB. Run migration scripts found in `sqlc/db-schema.sql` if needed.

**Database setup hints**
- Look at `sqlc/db-schema.sql` and `sqlc/db-queries.sql` to understand required tables.
- Use the DB connection settings in `connection/dev-config.json` for local dev.
- If `sqlc` is used, regenerate code after changing SQL: `sqlc generate` (if installed).
**Project API Flow — BriefBuletin_API**

This document explains how this specific project starts and how requests flow through the codebase. It contains concrete references to the actual files in this repo so a new developer (a fresher) can follow the wiring and configuration.

**Project Layout (where to look)**
- `main.go`: entry point and configuration loader (uses `viper` and accepts flags `-c`, `-v`, `-port`).
- `internal/service/`: business logic and HTTP handlers (examples: `routeService.go`, `rest_Service.go`, `commentService.go`).
- `internal/dbmodel/db_query`: generated/handwritten DB queries used by services.
- `internal/model`: typed request/response and auth claim structs.
- `util/`: helpers (DB connection wrapper, common utilities).
- `config.json` and `config/connection/dev-config.json`: environment configs. The app reads the config file given by `-c` (defaults to `./config.json`).

**High-level flow (exact for this repo)**
- `main.go` reads configuration (default: `./config.json`) and sets up `viper` with `AutomaticEnv()` and `SetEnvKeyReplacer`.
- `main.go` calls `service.NewAPIServer(configBytes, verbose)` to create the server instance and then calls `Serve(port)` to start the HTTP server.
- The server code wires a Gin router, registers endpoints from `internal/service/routeService.go` (and other service files), attaches middleware, and then listens on the configured port.

**Startup details (what `main.go` does)**
- Flags supported (see `main.go`):
  - `-c ./config.json` : config file path
  - `-v` : verbose logging (boolean)
  - `-port 7070` : port to run server on (defaults to 7070)
- `main.go` reads the file content and passes the raw JSON bytes into `service.NewAPIServer(configBytes, verbose)` so services can parse config values.

**Routes (where they are defined and how they map)**
- Public / open-api routes are added by `OpenAPIService.AddRouters` in `internal/service/routeService.go`. Examples:
  - `GET /api/category` -> `getAllCategory`
  - `GET /api/articles` -> `getApprovedArticles`
  - `POST /api/comment` -> `createComment`
  - `GET /api/all-comments` -> `getNewsComments`
- Admin / REST routes are added by `RESTService.AddRouters` in the same file. Examples:
  - `POST /api/auth/create` -> user creation
  - `POST /api/auth/login` -> login
  - `GET /api/active-comment` -> `approveComment` (requires auth/roles)

When you open `internal/service/routeService.go` you can see handlers call service methods and then `c.JSON(resp.StatusCode, resp)` — the handler returns an `APIResponse` typed struct with a status code and payload.

**Authentication: exact behavior from the code**
- Token creation: `RESTService.createJWTToken(userID, email, userName, role)` in `internal/service/authenticationService.go`:
  - Uses `s.jwtSigningKey` to sign a JWT (HMAC SHA256).
  - Uses shorter expiry for `ADMIN` (1 hour) and 24 hours for other roles.
  - Claims are built from `model.AuthorizationClaims` (includes `user_id`, `email`, `user_name`, `role`, and standard claims).
- Token validation / middleware helper: `RESTService.checkAuth(c *gin.Context)`:
  - First, it checks if the request URI is present in `s.bypassAuth`; if so, it returns `true` (allow).
  - If `s.jwtSigningKey == nil`, it returns `true` (the code treats missing signing key as allow-all).
  - Otherwise, it reads the `Authorization` header, expects `Bearer <token>`, parses token with `jwt.Parse` using the same `s.jwtSigningKey` and confirms `token.Valid`.
  - On successful parse it extracts claims and sets a typed user context under Gin key `"__USER_INFO__"` using `c.Set("__USER_INFO__", userInfo)`.
  - Helper `GetLoggedInUser(c)` reads that key and returns `*model.UserRoleInfo`.

This means protected service methods call `GetLoggedInUser` to access the current user's role and id. Example in `internal/service/commentService.go`: `user, _ := GetLoggedInUser(c)` is used to check `user.Role` before approving comments.

**Config keys and bypass list (from repo)**
- JWT signing key: config key `jwtKey` in `config.json` / `config/connection/dev-config.json`.
- Bypass list: config key `bypassAuth` is an array of URIs that `checkAuth` will allow without a token.
  - Example from `config/connection/dev-config.json` includes:
    - `/api/auth/create`, `/api/auth/login`, `/api/auth/resetpwd`, `/api/send-otp`, `/api/verify-otp`, `/api/verify-user`.
  - Example from root `config.json` includes the same auth endpoints plus API docs paths.

Important: the code also treats a missing signing key (`s.jwtSigningKey == nil`) as "allow all" — so if `jwtKey` is not loaded into the server, authentication is effectively disabled.

**Bypass / dev-mode behavior (exact)**
- The project uses the `bypassAuth` list from config. If a request URI exactly matches an entry in that slice, `RESTService.checkAuth` returns `true` and no JWT validation occurs.
- `jwtKey` missing behavior: if `s.jwtSigningKey` is not set, `checkAuth` returns `true` for all requests.

Security caution: both behaviors are convenient for local dev but must be disabled in production. Ensure `jwtKey` is set in production and `bypassAuth` contains only allowed public endpoints (or is empty).

**Service example: comment flow (exact code mapping)**
- POST `/api/comment` -> `OpenAPIService.createComment(c)` in `internal/service/commentService.go`:
  - Parses `model.CreateComment`, checks `ArticleID`, `UserName`, `UserEmail`.
  - Uses DB querier `auth.New(db)` (the generated DB wrapper) and calls `CreateCommentWithDefaults`.
  - Returns `BuildResponse200` or `BuildResponse500` accordingly.
- GET `/api/all-comments?article_id=<id>` -> `OpenAPIService.getNewsComments(c)`:
  - Calls `GetArticleDetails` then `GetApprovedCommentsByArticle` from DB layer.

**How to run the app locally (copy/paste commands for PowerShell)**
Make sure you are in the project root directory.

```pwsh
# Build binary
go build -o build\\exec\\briefbuletin main.go

# Run binary (uses default config path ./config.json; change with -c)
.\\build\\exec\\briefbuletin -c ./config.json -port 7070

# Or run directly with go (uses same flags)
go run main.go -c ./config.json -port 7070
```

If you want to use the dev config file: pass the `-c` flag with its path, e.g.: `-c ./config/connection/dev-config.json`.

**Where to look when things fail**
- `main.go` prints errors when it cannot read the config file — verify your `-c` path.
- If DB connection fails, check the `util/dbconnectionwraper.go` and the `dbhost`, `dbPort`, `dbuid`, `dbpassword` entries in your config file.
- If protected endpoints return `401`, check `jwtKey` in the config and ensure the client sends `Authorization: Bearer <token>`.
- If requests bypass auth unexpectedly, verify `bypassAuth` array and that `jwtKey` is present.

**Quick checklist for a new developer**
- Open `main.go` to see flags and how config bytes are passed to `service.NewAPIServer`.
- Open `internal/service/authenticationService.go` to understand `createJWTToken` and `checkAuth` behavior.
- Open `internal/service/routeService.go` to see all registered endpoints and which handlers implement them.
- Run the server locally with the dev config: `go run main.go -c ./config/connection/dev-config.json` and test the login endpoint to get a token.
- Use `GetLoggedInUser(c)` in debugger or code to inspect values injected by `checkAuth`.

If you want, I can now update this doc further by including small code snippets showing how `s.jwtSigningKey` is loaded in `NewAPIServer` (so we can confirm if it decodes the `jwtKey` value), or I can open `service.NewAPIServer` and paste the exact lines that perform config parsing. Which would you prefer?

---
File location: `README_API_FLOW.md` (project root). I updated it to reference the exact files and behavior found in this repository.