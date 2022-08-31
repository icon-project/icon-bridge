package errors

import "errors"

func New(err string) error {
	return errors.New("invalid block height")
}