# README.md

This file provides guidance to the reader when working with code in this repository.

## Overview

Package-locker system (low-level design exercise) in Go. Four flat files in `package main`: `locker.go`, `compartment.go`, `token.go`, `size.go`.

## Build / Run / Test

Current state: no `go.mod`, no `_test.go` files, not a git repo. `main.go` is a scenario runner printing PASS/FAIL lines.

- Run: `go run *.go` (no module, so `go run .` fails until `go mod init <name>`).
- `go vet ./...` and `gofmt -l .` work once a module exists (`token.go` struct fields are not gofmt-aligned).
- Tests, when added: `go test ./...`; single test: `go test -run TestName ./...`.

## Architecture

- `Locker` owns one map `index -> *compartmentSlot`, each slot pairing a `*Compartment` with its current `*AccessToken` (nil when empty). Index = position in the slice passed to `NewLocker`.
- `DepositPackage(size)`: first-fit scan by index for an unoccupied compartment of the **exact** size (no fallback to larger sizes), marks it occupied, stores a random 6-digit code with 7-day expiry in the slot.
- `Pickup(code)`: linear scan of slots for the code, checks expiry, opens compartment, then `clearDeposit` frees the compartment and nils the token.
- `OpenExpiredCompartments()`: opens compartments with expired tokens but does **not** free them or clear tokens.
- `AccessToken` still has a `compartment` field (from `token.go`), but `Locker` passes `nil` and never uses it; the slot map is the source of truth.
- `Compartment.Open()` is an empty stub (hardware hook).

Known gaps worth keeping in mind when changing behavior: token codes are not checked for collisions against existing map entries; no concurrency safety (plain map, no mutex); expiry uses `time.Now()` directly, so tests need a clock abstraction or short expirations.
