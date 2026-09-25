package main

import (
	"fmt"
	"time"
)

func main() {
	locker := NewLocker([]*Compartment{
		NewCompartment(Small),
		NewCompartment(Medium),
		NewCompartment(Large),
		NewCompartment(Small),
	})

	// Deposit fills compartments of matching size in index order.
	code1, err := locker.DepositPackage(Small)
	check("deposit small #1", err == nil, err)
	code2, err := locker.DepositPackage(Small)
	check("deposit small #2", err == nil, err)
	_, err = locker.DepositPackage(Small)
	check("deposit small #3 fails (none left)", err != nil, err)
	check("slot 0 occupied", locker.slots[0].compartment.IsOccupied(), nil)
	check("slot 3 occupied", locker.slots[3].compartment.IsOccupied(), nil)

	// Pickup frees the compartment and invalidates the code.
	err = locker.Pickup(code1)
	check("pickup code1", err == nil, err)
	check("slot 0 freed", !locker.slots[0].compartment.IsOccupied() && locker.slots[0].accessToken == nil, nil)
	err = locker.Pickup(code1)
	check("pickup code1 again fails", err != nil, err)

	// Invalid codes are rejected.
	check("pickup empty code fails", locker.Pickup("") != nil, nil)
	check("pickup unknown code fails", locker.Pickup("not-a-code") != nil, nil)

	// Freed compartment is reused.
	_, err = locker.DepositPackage(Small)
	check("deposit small reuses slot 0", err == nil && locker.slots[0].compartment.IsOccupied(), err)

	// Expired token: pickup fails, compartment stays occupied.
	locker.slots[3].accessToken.expiration = time.Now().Add(-time.Minute)
	err = locker.Pickup(code2)
	check("pickup expired code fails", err != nil, err)
	check("expired slot still occupied", locker.slots[3].compartment.IsOccupied(), nil)
	locker.OpenExpiredCompartments()

	// Other sizes work independently.
	code3, err := locker.DepositPackage(Large)
	check("deposit large", err == nil, err)
	check("pickup large", locker.Pickup(code3) == nil, nil)
}

func check(name string, ok bool, err error) {
	status := "PASS"
	if !ok {
		status = "FAIL"
	}
	if err != nil {
		fmt.Printf("[%s] %s (err: %v)\n", status, name, err)
		return
	}
	fmt.Printf("[%s] %s\n", status, name)
}
