package db

import "time"

func (i Invitation) IsPending() bool {
	if i.AcceptedAt != nil || i.DeclinedAt != nil || i.CanceledAt != nil {
		return false
	}
	if time.Now().Before(i.ExpiresAt) {
		return false
	}
	return true
}
