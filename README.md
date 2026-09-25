# Locker

Package-locker system (low-level design exercise) in Go. Flat `package main`: `locker.go`, `compartment.go`, `token.go`, `size.go`, plus `main.go` as a scenario runner.

## Run

```sh
go run .
```

`main.go` exercises the locker and prints `[PASS]` / `[FAIL]` lines. Most scenarios are currently commented out; uncomment the block in `main()` to run them all (it also needs `"time"` re-added to the imports).

There are no `_test.go` files yet. Once added: `go test ./...`, or a single test with `go test -run TestName ./...`.

## Usage

```go
locker := NewLockerWithSizes(Small, Medium, Large, Small)

code, err := locker.DepositPackage(Small) // returns 6-digit pickup code
err = locker.Pickup(code)                 // frees the compartment
```

`NewLocker([]*Compartment)` is still available when you want to build compartments yourself.

## Architecture

- `Locker` owns one map `index -> *compartmentSlot`, each slot pairing a `*Compartment` with its current `*AccessToken` (nil when empty). Index = position in the list given to the constructor.
- `DepositPackage(size)`: first-fit scan by index for an unoccupied compartment of the **exact** size (no fallback to larger sizes), marks it occupied, stores a random 6-digit code with 7-day expiry in the slot.
- `Pickup(code)`: linear scan of slots for the code, checks expiry, opens compartment, then frees it and clears the token.
- `OpenExpiredCompartments()`: opens compartments with expired tokens but does **not** free them or clear tokens.
- `AccessToken` still has a `compartment` field, but `Locker` passes `nil` and never uses it; the slot map is the source of truth.
- `Compartment.Open()` is an empty stub (hardware hook).

## Known gaps

- Token codes are not checked for collisions with codes already in use.
- No concurrency safety (plain map, no mutex).
- Expiry uses `time.Now()` directly, so tests need a clock abstraction or must edit the expiration.
