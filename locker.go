package main

import (
	"fmt"
	"math/rand"
	"time"
)

type compartmentSlot struct {
	compartment *Compartment
	accessToken *AccessToken
}

type Locker struct {
	slots map[int]*compartmentSlot
}

func NewLocker(compartments []*Compartment) *Locker {
	slots := make(map[int]*compartmentSlot, len(compartments))
	for i, c := range compartments {
		slots[i] = &compartmentSlot{compartment: c}
	}
	return &Locker{
		slots: slots,
	}
}

func NewLockerWithSizes(sizes ...Size) *Locker {
	compartments := make([]*Compartment, len(sizes))
	for i, size := range sizes {
		compartments[i] = NewCompartment(size)
	}
	return NewLocker(compartments)
}

func (l *Locker) DepositPackage(size Size) (string, error) {
	slot := l.getAvailableSlot(size)
	if slot == nil {
		return "", fmt.Errorf("no available compartment of size %s", size)
	}

	slot.compartment.Open()
	slot.compartment.MarkOccupied()
	slot.accessToken = l.generateAccessToken()

	return slot.accessToken.GetCode(), nil
}

func (l *Locker) Pickup(tokenCode string) error {
	if tokenCode == "" {
		return fmt.Errorf("invalid access token code")
	}

	slot := l.findSlotByCode(tokenCode)
	if slot == nil {
		return fmt.Errorf("invalid access token code")
	}

	if slot.accessToken.IsExpired() {
		return fmt.Errorf("access token has expired")
	}

	slot.compartment.Open()
	l.clearDeposit(slot)
	return nil
}

func (l *Locker) OpenExpiredCompartments() {
	for _, slot := range l.slots {
		if slot.accessToken != nil && slot.accessToken.IsExpired() {
			slot.compartment.Open()
		}
	}
}

func (l *Locker) getAvailableSlot(size Size) *compartmentSlot {
	for i := 0; i < len(l.slots); i++ {
		slot := l.slots[i]
		if slot.compartment.GetSize() == size && !slot.compartment.IsOccupied() {
			return slot
		}
	}
	return nil
}

func (l *Locker) findSlotByCode(tokenCode string) *compartmentSlot {
	for _, slot := range l.slots {
		if slot.accessToken != nil && slot.accessToken.GetCode() == tokenCode {
			return slot
		}
	}
	return nil
}

func (l *Locker) generateAccessToken() *AccessToken {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiration := time.Now().Add(7 * 24 * time.Hour)
	return NewAccessToken(code, expiration, nil)
}

func (l *Locker) clearDeposit(slot *compartmentSlot) {
	slot.compartment.MarkFree()
	slot.accessToken = nil
}
