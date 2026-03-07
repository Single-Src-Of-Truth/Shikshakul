# Shikshakul Core Library (core-go)

This module contains shared utilities, middleware, and core domain logic utilized across all Shikshakul Go microservices (IAM, Academics, Utility). It is integrated into the monorepo using Go Workspaces (go.work).

## API Versioning Engine

Shikshakul utilizes a global, date-based API versioning strategy (e.g., `Skl-Service-Version: 2026-03-06-a1b2c`). This allows client applications to lock into specific API contracts while the backend safely evolves.

The versioning system is fully in-memory, resulting in zero-latency lookups. It is managed via an internal CLI tool that generates compiled Go code from a JSON source of truth.

### Important Guidelines

- Do not manually edit `pkg/versioning/registry.go`. It is an auto-generated artifact.
- Do not manually edit `pkg/versioning/versions.json`. Let the CLI tool manage the state and formatting to prevent schema corruption.

## Managing Versions via CLI

The versioning CLI tool is located at `lib/core-go/tools/versioncli`. Run these commands from the root of the monorepo or from within the versioncli directory.

### 1. Create a New Release

When introducing breaking changes or significant platform updates, generate a new global API version. This automatically sets the new version to ACTIVE, marks it as the system default, and recompiles the registry.

```bash
go run lib/core-go/tools/versioncli/main.go release --desc "Brief description of the release changes"
```

### 2. Deprecate an Older Release

To mark an older version as deprecated (which can later be used to trigger warnings for clients still using it), use the deprecate command with the target version ID.

```bash
go run lib/core-go/tools/versioncli/main.go deprecate --id "2026-03-06-a1b2c"
```

### 3. Manual Code Generation

If you pull the repository and need to rebuild the Go registry from the JSON state file without creating a new release, run the generate command.

```bash
go run lib/core-go/tools/versioncli/main.go generate
```

Alternatively, from within the `lib/core-go/versioning` directory, you can run standard Go generation:

```bash
go generate
```

## Middleware Integration

To enforce versioning on a microservice, import the versioning package and attach the middleware to the global Gin router.

```go
import "github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/versioning"

func main() {
    router := gin.Default()

    // Apply global version enforcement
    router.Use(versioning.EnforceVersion())

    // Example endpoint to expose version history to frontend clients
    router.GET("/versions", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "default": versioning.DefaultVersion,
            "history": versioning.Registry,
        })
    })

    // ... route definitions
}
```

The middleware handles the following automatically:

- Intercepts the `Skl-Service-Version` header.
- Rejects requests with invalid or non-existent versions (400 Bad Request).
- Falls back to the compiled `DefaultVersion` if the client provides no header.
- Injects the resolved version into the response headers.
- Sets the `api_version` key in the Gin Context for downstream controller logic.
