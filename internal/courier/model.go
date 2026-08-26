package courier

import "uuid"

type provider string

const (
	GLSSIProvider provider = "gls_si"
)

type Courier struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Provider       string
	// Name           string
	AccountNumber string
	APIKey        string
}
