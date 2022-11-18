package errors

import "errors"

func New(err string) error {
	return errors.New(err)
}

var (
	ErrUnknownTransaction = errors.New("UNKNOWN_TRANSACTION")
	ErrUnknownBlock       = errors.New("UNKNOWN_BLOCK")
	ErrUnknown            = errors.New("UNKNOWN")
)

func Is(err, target error) bool {
	return err.Error() == target.Error()
}
