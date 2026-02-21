# Cinema.log AI Agent Instructions

## Architecture Overview

This is a **full-stack film rating application** with:

- **Frontend**: Angular 21 with PrimeNG (standalone components, signals-based state)
- **Backend**: Go 1.24+ REST API following vertical slice architecture
- **Database**: PostgreSQL managed via Docker Compose with Goose migrations
- **Auth**: GitHub OAuth with JWT tokens (access + refresh)

## Development Workflow

### Starting the dev environment

Use the **root-level script** (not individual `ng serve` or `go run`):

```bash
./run-dev.zsh  # Starts Angular, Docker Postgres, and Go server in parallel
```

The script:

1. Launches Angular dev server in a new Terminal window (`npm run start` at port 4200)
2. Opens Docker Desktop and spins up `psql_bp` container
3. Runs migrations automatically via `database.NewWithMigrations()`
4. Starts Go API at port 8080

### Stopping the dev environment

```bash
./stop-dev.zsh  # Kills ports 4200/8080, stops Docker containers
```

### Running tests

```bash
# Backend - use Makefile or direct go test
cd cinema.log.server.golang
make test  # OR: go test ./... -v

# Frontend - Unit/Component tests
cd cinema.log.client
npm run test  # Uses vitest with Angular TestBed for component/service tests

# Frontend - E2E tests with Playwright
cd cinema.log.client
npx playwright test              # Run all e2e tests
npx playwright test --ui         # Interactive UI mode
npx playwright test --headed     # Run with browser visible
npx playwright show-report       # View HTML test report
```

## Backend Architecture (Go)

### Vertical Slice Pattern

Each feature domain is self-contained in `internal/{feature}/`:

```
internal/films/
  ├── film_handler.go       # HTTP layer (request/response)
  ├── film_service.go       # Business logic
  ├── film_store.go         # Database queries
  ├── film_handler_test.go  # Unit tests with mocks
  └── film_service_test.go
```

### Dependency Injection (DI)

All wiring happens in [internal/server/server.go](../cinema.log.server.golang/internal/server/server.go):

```go
// Pattern: DB -> Store -> Service -> Handler
filmStore := films.NewStore(db)
filmService := films.NewService(filmStore, graphService)
filmHandler := films.NewHandler(filmService, ratingService)
```

**Cross-slice dependencies** use interface types defined in the _consuming_ handler/service (e.g., `FilmService` interface in [review_handler.go](../cinema.log.server.golang/internal/reviews/review_handler.go)).

### Adding New Routes

1. Create handler method in `{feature}/handler.go`
2. Register in [internal/server/routes.go](../cinema.log.server.golang/internal/server/routes.go):
   ```go
   mux.HandleFunc("GET /new-route", s.newHandler.NewMethod)
   ```
3. Add to `isAuthExempt()` if no auth required
4. **Batch endpoints**: Use descriptive names like `/ratings/compare-films-batch` for bulk operations

### Testing Strategy

- **Unit tests**: Use hand-written mock structs (e.g., `mockFilmService` in [film_handler_test.go](../cinema.log.server.golang/internal/films/film_handler_test.go)) with function fields for custom behavior
- **Integration tests**: Use testcontainers ([utils/test_utils.go](../cinema.log.server.golang/internal/utils/test_utils.go)) - see `StartTestPostgres()` for real DB setup
- **Mock Updates**: When adding new interface methods, ensure mocks in test files implement them:
  - Handler tests: Update `mock*Service` structs
  - Service tests: Update `mock*Store` structs with function fields and corresponding methods
- **Frontend unit tests**: Vitest with Angular TestBed for components/services
- **E2E tests**: Playwright for full user journey testing (see E2E Testing section below)

### Database Migrations

- Location: `internal/migration/goose/*.sql`
- Auto-applied on startup via `database.NewWithMigrations()`
- Manual: Use `database.New()` and run Goose CLI separately

### Auth Middleware

[routes.go](../cinema.log.server.golang/internal/server/routes.go) validates JWT on all routes except:

- `/auth/github-login`, `/auth/github-callback`, `/auth/refresh-token`
- `/auth/dev/login` (dev-only bypass)

Authenticated user ID available via:

```go
userId, err := middleware.GetUserIDFromContext(r.Context())
```

## Frontend Architecture (Angular)

### Modern Angular Patterns

- **Standalone components** (no NgModules except legacy)
- **Inject() for DI** (no constructor injection in services/components)
- **Signals** for reactive state (`signal()`, `computed()`)
  ```typescript
  currentUser = signal<User | null>(null); // in AuthService
  isLoading = signal(true); // in components
  ```
- **Functional guards** ([auth.guard.ts](../cinema.log.client/src/app/guards/auth.guard.ts)) - checks `currentUser()` signal or calls `/auth/me`

### Services & HTTP

- All HTTP calls use `withCredentials: true` for cookies (JWT storage)
- Error handling via [utils/error-handler.util.ts](../cinema.log.client/src/app/utils/error-handler.util.ts):
  - `handleHttpError()` - unexpected errors (show user-friendly message)
  - `handleExpectedError()` - auth failures (silent, e.g., 401 on `/auth/me`)
- Base API URL from `import.meta.env.NG_APP_API_URL`

### Component Structure

```
components/{feature}/
  ├── {feature}.component.ts     # Component logic
  ├── {feature}.component.html   # Template (PrimeNG components)
  ├── {feature}.component.scss   # Styles (uses _variables.scss, _mixins.scss)
  └── {feature}.component.spec.ts
```

### PrimeNG Integration

- Theme: Aura (dark mode via `.dark-mode` class)
- Import individual components: `import { ButtonModule } from 'primeng/button';`
- Config in [app.config.ts](../cinema.log.client/src/app/app.config.ts) with ripple effects enabled

## Environment Variables

Backend requires (`.env` in `cinema.log.server.golang/`):

```env
PORT=8080
TOKEN_SECRET=your_jwt_secret
FRONTEND_URL=http://localhost:4200
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=test_database
BLUEPRINT_DB_USERNAME=test_user
BLUEPRINT_DB_PASSWORD=test_password
BLUEPRINT_DB_SCHEMA=public
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
TMDB_ACCESS_TOKEN=xxx  # For film search API
```

Frontend: .env file with:

```env
NG_APP_API_URL=http://localhost:8080
```

## Key Integration Points

### Auth Flow

1. User clicks "Login" → redirects to `/auth/github-login`
2. Backend handles OAuth callback → sets JWT cookie
3. Frontend guards call `/auth/me` to validate session
4. `AuthService.currentUser` signal stores user globally

### Film Rating System

- **Initial rating**: Created when user submits review via `POST /reviews` (triggers graph update)
- **Bulk Comparisons**: `POST /ratings/compare-films-batch` - batch processing of film comparisons with:
  - Strict sequential ELO calculation (preserves K-factor progression)
  - All-or-nothing transaction semantics
  - Automatic duplicate filtering
  - 50-film cap per batch
- **Graph**: Film relationships stored in `film_graph` table, visualized with vis-network

#### Bulk Review Feature (Review Component)

- **Default Mode**: Bulk comparison mode with toggle to sequential
- **UI Components**: Sticky target film at top, progress counter (e.g., "3 / 50 films loaded"), comparison cards with Better/Same/Worse buttons
- **Progressive Loading**: Initial 10 films, expandable to 50 via "Load More" button
- **State Management**: Uses Angular signals (`isBulkMode`, `bulkSelections`)
- **Persistence**: Mode preference stored in localStorage (`comparisonMode`)
- **Partial Submissions**: Users can submit any number of completed comparisons (validation allows partial selections)
- **Empty State**: Auto-redirects to profile when 0 films available for comparison

### External API (TMDB)

Film search via `GET /films/search?f=query` proxied through backend [films/film_service.go](../cinema.log.server.golang/internal/films/film_service.go) to avoid exposing API keys
E2E Testing (Playwright)

### Test Structure

Tests are organized by feature in `cinema.log.client/tests/`:

```
tests/
  ├── auth.setup.ts              # Global auth setup (runs before all tests)
  ├── seed.spec.ts               # Verifies auth is working
  ├── authentication/            # Auth-related tests
  │   ├── protected-routes.spec.ts
  │   └── sign-out.spec.ts
  ├── film-search/               # Film search feature tests
  │   ├── search-by-title.spec.ts
  │   ├── select-film.spec.ts
  │   └── no-results.spec.ts
  ├── review/                    # Review submission tests
  │   ├── submit-review.spec.ts
  │   ├── star-rating-selection.spec.ts
  │   └── validation-no-rating.spec.ts
  ├── comparison-flow/           # Bulk comparison tests
  │   └── submit-then-bulk.spec.ts
  ├── profile/                   # User profile tests
  │   └── profile-information.spec.ts
  ├── navigation/                # Navigation tests
  │   └── logo-home.spec.ts
  └── utils/                     # Test utilities
      └── test-helpers.ts
```

### Authentication Setup

Tests use a **global authentication setup** via [auth.setup.ts](../cinema.log.client/tests/auth.setup.ts):

- Runs once before all tests (defined in `playwright.config.ts` as `setup` project)
- Uses `/auth/dev/login` endpoint to bypass OAuth (creates test user automatically)
- Saves authentication state to `.auth/user.json`
- All test projects depend on the `setup` project and load stored auth state

```typescript
// Tests automatically have authenticated context
test("some test", async ({ page }) => {
  await page.goto("/"); // Already logged in!
});
```

### Test Helpers & Utilities

[tests/utils/test-helpers.ts](../cinema.log.client/tests/utils/test-helpers.ts) provides:

- **`ensureFilmExists(page, filmTitle, rating?, reviewText?)`**: Idempotent test data creation
  - Searches for film and adds it to user's collection if not present
  - Makes tests independent and repeatable
  - Example:
    ```typescript
    await ensureFilmExists(page, "The Dark Knight", 5, "Great film!");
    ```

### Configuration

[playwright.config.ts](../cinema.log.client/playwright.config.ts) settings:

- **Base URL**: `http://localhost:4200` (or `process.env.BASE_URL`)
- **Test isolation**: Fully parallel execution
- **Projects**: Setup project + Chromium (with stored auth state)
- **CI Configuration**:
  - 3 retries on CI, 0 locally
  - Sequential execution on CI (`workers: 1`)
  - Fails build if `test.only` found
- **Reporting**: HTML reports (`npx playwright show-report`)

### Test Patterns

**Test structure**:

```typescript
import { test, expect } from "@playwright/test";

test.describe("Feature Name", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("user can perform action", async ({ page }) => {
    // Use data-testid attributes for reliable selectors
    await page.getByTestId("navbar-search-link").click();
    await expect(page.getByRole("heading", { name: "Search" })).toBeVisible();
  });
});
```

**Selector strategy**:

- **Preferred**: `getByTestId()` for custom `data-testid` attributes
- **Fallback**: `getByRole()` for semantic elements (headings, buttons, etc.)
- **Avoid**: CSS class selectors (too brittle)

**Common patterns**:

- Navigation tests verify routing and redirects
- Form tests check validation and submission
- Protected route tests verify auth guards
- Use `waitForURL()` for navigation assertions
- Use `ensureFilmExists()` for idempotent test data

### Running E2E Tests

```bash
# Ensure dev environment is running first
./run-dev.zsh

# In another terminal:
cd cinema.log.client
npx playwright test                    # Run all tests
npx playwright test --ui               # Interactive mode with time-travel debugging
npx playwright test film-search/       # Run specific directory
npx playwright test --grep "search"    # Run tests matching pattern
npx playwright test --headed           # See browser during test
npx playwright show-report             # View last test report
```

##

## Common Patterns

### Backend

- **Error handling**: Return `fmt.Errorf("description: %w", err)` for wrapping
- **Context propagation**: Always pass `ctx context.Context` as first param
- **SQL queries**: Use `pgx` driver with prepared statements (`$1`, `$2` placeholders)
- **HTTP responses**: Use utility from [utils/json.go](../cinema.log.server.golang/internal/utils/json.go)
- **Transactions**: Use `BeginTx()` for multi-step database operations (e.g., batch updates)

### Frontend

- **HTTP error handling**: Always use `catchError()` with error handler utilities
- **Loading states**: Use signals for `isLoading`, `errorMessage`
- **Router navigation**: Inject `Router` for programmatic nav (e.g., after login success)
- **Form handling**: Angular signals for reactive form state (no template-driven forms)
- **Batch Operations**: Submit via dedicated bulk endpoints for better performance and transaction integrity

## Migration/Deployment Notes

- Production backend: Railway deployment
- Frontend build: `ng build` → outputs to `dist/`
- DB migrations: Goose applies in order based on timestamp prefix
- Docker: Use `docker compose` (v2 syntax) in scripts, not `docker-compose`
