package utils

import (
	"errors"

	"github.com/lib/pq"
)

const (
	ErrUniqueViolation = "unique_violation"
)

func IsErrorCode(err error, errcode string) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		println(pgErr.Code)
		return pgErr.Code.Name() == errcode
	}
	return false
}
