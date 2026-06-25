package service

import (
	"errors"

	"gorm.io/gorm"
)

// notFoundOr returns a friendly "not found" error when err is gorm's
// ErrRecordNotFound, otherwise returns the original error — so real DB failures
// (timeouts, connection drops) are not masked as "record not found".
func notFoundOr(err error, notFoundMsg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(notFoundMsg)
	}
	return err
}
