# leakcheck

[![CI](https://github.com/headlesslab/leakcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/headlesslab/leakcheck/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/headlesslab/leakcheck.svg)](https://pkg.go.dev/github.com/headlesslab/leakcheck)

A library for monitoring the runtime goroutine stack: wait for goroutines to exit, detect leaks in tests.

- `context.Context` first design
- concurrent leak detection
- no dependencies
- handy low-level APIs to extend the library

## Provenance

`leakcheck` is a Satellite module of [wand](https://github.com/headlesslab/wand). It is a snapshot of [gotrace](https://github.com/ysmood/gotrace) `v0.6.0` (commit `b96e976`), imported without git history and maintained here by headlesslab; see [NOTICE](NOTICE). The module path changed, the exported identifiers did not.

Leak detection attributes goroutines to the test that started them through the runtime's ancestor tracebacks, so the test binary must run with `GODEBUG=tracebackancestors=N` (N large enough, such as 1000); `IgnoreNonChildren` panics otherwise.

## Usage

```go
import "github.com/headlesslab/leakcheck"
```

The tests are the documentation: see [examples_test.go](examples_test.go).

## Versioning

`v0.x`: no breaking change without a major version bump. Releases are cut by hand when the module changes; wand pins each release exactly.

## License

MIT, see [LICENSE](LICENSE).
