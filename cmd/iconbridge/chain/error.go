package chain

import (
	"errors"
	"strings"
)

var (
	// Common errors
	ErrInsufficientBalance   = errors.New("InsufficientBalance")
	ErrGasLimitExceeded      = errors.New("GasLimitExceeded")
	ErrBlockGasLimitExceeded = errors.New("BlockGasLimitExceeded")

	// BMC errors
	ErrBMCRevertLastOwner                 = errors.New("LastOwner")
	ErrBMCRevertUnauthroized              = errors.New("Unauthorized")
	ErrBMCRevertInvalidAddress            = errors.New("InvalidAddress")
	ErrBMCRevertNotExistsPermission       = errors.New("NotExistsPermission")
	ErrBMCRevertAlreadyExistsBSH          = errors.New("AlreadyExistsBSH")
	ErrBMCRevertNotExistsBSH              = errors.New("NotExistsBSH")
	ErrBMCRevertAlreadyExistsLink         = errors.New("AlreadyExistsLink")
	ErrBMCRevertNotExistsLink             = errors.New("NotExistsLink")
	ErrBMCRevertInvalidParam              = errors.New("InvalidParam")
	ErrBMCRevertAlreadyExistRoute         = errors.New("AlreadyExistRoute")
	ErrBMCRevertNotExistRoute             = errors.New("NotExistRoute")
	ErrBMCRevertInvalidSn                 = errors.New("InvalidSn")
	ErrBMCRevertParseFailure              = errors.New("ParseFailure")
	ErrBMCRevertInvalidRxHeight           = errors.New("InvalidRxHeight")
	ErrBMCRevertInvalidSeqNumber          = errors.New("InvalidSeqNumber")
	ErrBMCRevertNotExistsInternalHandler  = errors.New("NotExistsInternalHandler")
	ErrBMCRevertAlreadyExistsBMCPeriphery = errors.New("AlreadyExistsBMCPeriphery")
	ErrBMCRevertUnknownHandleBTPError     = errors.New("UnknownHandleBTPError")
	ErrBMCRevertUnknownHandleBTPMessage   = errors.New("UnknownHandleBTPMessage")
	ErrBMCRevertUnreachable               = errors.New("Unreachable:")
)

func RevertError(msg string) error {
	for _, err := range []error{
		ErrBMCRevertLastOwner,
		ErrBMCRevertUnauthroized,
		ErrBMCRevertInvalidAddress,
		ErrBMCRevertNotExistsPermission,
		ErrBMCRevertAlreadyExistsBSH,
		ErrBMCRevertNotExistsBSH,
		ErrBMCRevertAlreadyExistsLink,
		ErrBMCRevertNotExistsLink,
		ErrBMCRevertInvalidParam,
		ErrBMCRevertAlreadyExistRoute,
		ErrBMCRevertNotExistRoute,
		ErrBMCRevertInvalidSn,
		ErrBMCRevertParseFailure,
		ErrBMCRevertInvalidRxHeight,
		ErrBMCRevertInvalidSeqNumber,
		ErrBMCRevertNotExistsInternalHandler,
		ErrBMCRevertAlreadyExistsBMCPeriphery,
		ErrBMCRevertUnknownHandleBTPError,
		ErrBMCRevertUnknownHandleBTPMessage,
		ErrBMCRevertUnreachable,
	} {
		if strings.HasPrefix(msg, err.Error()) {
			return err
		}
	}
	return nil
}
