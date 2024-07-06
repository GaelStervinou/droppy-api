package errors2

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolationErr = "23505"
)

func IsErrorCode(err error, errcode string) bool {
	var pgerr *pgconn.PgError
	errors.As(err, &pgerr)

	return errcode == pgerr.Code
}
