package errors

import "fmt"

type BlockNotificationError struct {
	Offset int

	Err error
}

func (e *BlockNotificationError) Error() string {
	return fmt.Sprintf("offset %d: err %v", e.Offset, e.Err)
}
