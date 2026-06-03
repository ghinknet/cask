# cask

Cache, abstracted simply.

`cask` is a small, driver-based caching abstraction for Go. It defines a tiny
set of interfaces for cache operations and lets concrete backends (Redis,
in-memory, etc.) plug in as **drivers**. Your application talks to a single
`Namespace` value while the actual storage backend is resolved automatically
from the client you hand in.

## Features

- **Driver-based** – backends register themselves and are matched to a client
  by type, so swapping backends needs no code changes at the call site.
- **Namespaces** – build hierarchical, colon-joined keys (`root:child:leaf`)
  for multi-tenant or multi-feature key isolation.
- **Minimal surface** – the core API is just `Get` / `Set` / `Del` / `Exists`,
  with optional `Expire` / `TTL` for backends that support expiration.
- **Pluggable TTL semantics** – durations use [`timex`](https://go.gh.ink/timex),
  which can express finite, zero, and positive/negative-infinite durations.

## Installation

```sh
go get go.gh.ink/cask
```

To actually store anything you also need a driver, e.g. the Redis driver:

```sh
go get go.gh.ink/cask/redis
```

## Quick start

```go
package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.gh.ink/cask"
	_ "go.gh.ink/cask/redis" // registers the "redis" driver via init()
	"go.gh.ink/timex"
)

func main() {
	// Pass any supported client; cask picks the matching driver automatically.
	root, err := cask.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), "app")
	if err != nil {
		panic(err) // errors.ErrDriverNotRegistered if no driver matches the client
	}

	// Derive a child namespace -> key "app:session"
	ns := root.Namespace("session")

	ctx := context.Background()
	_ = ns.Set(ctx, []byte("hello"), timex.NewPosInfDuration())
	val, _ := ns.Get(ctx)
	fmt.Println(string(val)) // hello
}
```

## Core concepts

### Namespace

A `Namespace` is the value your code interacts with. It bundles three things:

- the **raw client** you passed to `cask.New` (retrievable via `Raw()`),
- a **key path** — a `[]string` joined with `:` and exposed via `Key()`,
- the resolved **`Adapter`** (embedded), which provides the actual cache methods.

`New(client, key...)` walks every registered driver and asks it to build an
adapter for the given client. The first driver that accepts the client wins; if
none match, `errors.ErrDriverNotRegistered` is returned.

`Namespace(key...)` returns a **child** namespace. It is copy-on-write: the
parent is never mutated, the key path is extended, and the adapter is
re-resolved against the new key. This makes it safe to fan out namespaces from a
shared root.

```go
root, _ := cask.New(client, "app")        // key: "app"
users := root.Namespace("users")           // key: "app:users"
alice := users.Namespace("alice")          // key: "app:users:alice"
fmt.Println(alice.Key())                    // app:users:alice
```

## Package layout

| Package            | Responsibility                                                                 |
| ------------------ | ------------------------------------------------------------------------------ |
| `cask` (root)      | `New` constructor, `Namespace` type, and namespace chaining (`Namespace`, `Key`, `Raw`). |
| `model`            | The contract interfaces: `BaseStore`, `Expirer`, `Adapter`, `Driver`, `NamespaceInfo`. |
| `driver`           | The driver registry: `Register(name, driver)` and `List()`.                    |
| `errors`           | Shared sentinel errors, currently `ErrDriverNotRegistered`.                    |
| `internal/state`   | Process-wide `Drivers` map backing the registry (internal, not importable).    |
| `example`          | A reference driver showing the minimal driver/adapter implementation.          |

## Interfaces

Defined in the `model` package:

```go
type BaseStore interface {
	Get(ctx context.Context) ([]byte, error)
	Set(ctx context.Context, value []byte, ttl timex.Duration) error
	Del(ctx context.Context) (bool, error)
	Exists(ctx context.Context) (bool, error)
}

type Expirer interface {
	Expire(ctx context.Context, ttl timex.Duration) error
	TTL(ctx context.Context) (timex.Duration, error)
}

type Adapter interface {
	BaseStore
}

// NamespaceInfo is what a driver sees: the namespaced key string.
type NamespaceInfo interface {
	Key() string
}

// Driver builds an Adapter for a given client, or reports that it can't.
type Driver interface {
	NewAdapter(client any, ns NamespaceInfo) (adapter Adapter, ok bool)
}
```

`Adapter` requires `BaseStore`. `Expirer` is optional capability: backends that
support TTL management implement it, and callers reach it via a type assertion,
e.g. `ns.Adapter.(model.Expirer)`.

## Writing a driver

A driver is two pieces — an `Adapter` that performs operations and a `Driver`
that matches a client type and constructs the adapter — plus a one-line
registration. See the `example` package for the canonical skeleton:

```go
// 1. Adapter implements model.Adapter (BaseStore).
type Adapter struct{}

func (a Adapter) Get(ctx context.Context) ([]byte, error)                   { /* ... */ }
func (a Adapter) Set(ctx context.Context, v []byte, ttl timex.Duration) error { /* ... */ }
func (a Adapter) Del(ctx context.Context) (bool, error)                     { /* ... */ }
func (a Adapter) Exists(ctx context.Context) (bool, error)                  { /* ... */ }

// 2. Driver matches a client type and returns ok=true when it can handle it.
type Driver struct{}

func (d Driver) NewAdapter(client any, ns model.NamespaceInfo) (model.Adapter, bool) {
	if c, ok := client.(MyClientType); ok {
		return Adapter{ /* keep c and ns */ }, true
	}
	return nil, false
}

// 3. Register it. Drivers commonly do this in init() so a blank import is enough:
func init() { driver.Register("myname", Driver{}) }
```

> The bundled `example` driver registers via an exported `Trigger()` function
> instead of `init()`, so it is only active when explicitly invoked — useful as
> a reference without polluting the global registry. Real drivers (like
> `cask-redis`) register in `init()`.

## Errors

| Error                          | Meaning                                                       |
| ------------------------------ | ------------------------------------------------------------- |
| `errors.ErrDriverNotRegistered` | `New` found no registered driver willing to handle the client. |

## License

Apache License 2.0. See [LICENSE](LICENSE).
