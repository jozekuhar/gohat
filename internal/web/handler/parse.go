package handler

import (
	"net/http"
	"uuid"
)

func pathValueUUID(r *http.Request, key string) (uuid.UUID, error) {
	valueStr := r.PathValue(key)
	valueUUID, err := uuid.Parse(valueStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return valueUUID, nil
}
