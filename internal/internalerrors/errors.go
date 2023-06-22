package internalerrors

import (
	"errors"

	"gorm.io/gorm"
)

var ErrorInternal = errors.New("internal server error")

func ProcessErrorToReturn(err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorInternal
	}

	return err
}
