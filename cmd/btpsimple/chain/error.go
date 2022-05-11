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
	ErrBMCRevertLastOwner                 = errors.New("BMCRevertLastOwner")
	ErrBMCRevertUnauthroized              = errors.New("BMCRevertUnauthorized")
	ErrBMCRevertInvalidAddress            = errors.New("BMCRevertInvalidAddress")
	ErrBMCRevertNotExistsPermission       = errors.New("BMCRevertNotExistsPermission")
	ErrBMCRevertAlreadyExistsBSH          = errors.New("BMCRevertAlreadyExistsBSH")
	ErrBMCRevertNotExistsBSH              = errors.New("BMCRevertNotExistsBSH")
	ErrBMCRevertAlreadyExistsLink         = errors.New("BMCRevertAlreadyExistsLink")
	ErrBMCRevertNotExistsLink             = errors.New("BMCRevertNotExistsLink")
	ErrBMCRevertInvalidParam              = errors.New("BMCRevertInvalidParam")
	ErrBTPRevertAlreadyExistRoute         = errors.New("BTPRevertAlreadyExistRoute")
	ErrBTPRevertNotExistRoute             = errors.New("BTPRevertNotExistRoute")
	ErrBMCRevertInvalidSN                 = errors.New("BMCRevertInvalidSN")
	ErrBMCRevertParseFailure              = errors.New("BMCRevertParseFailure")
	ErrBMCRevertRxSeqLowerThanExpected    = errors.New("BMCRevertRxSeqLowerThanExpected")
	ErrBMCRevertRxSeqHigherThanExpected   = errors.New("BMCRevertRxSeqHigherThanExpected")
	ErrBMCRevertInvalidRxHeight           = errors.New("BMCRevertRxHeightLowerThanExpected")
	ErrBMCRevertNotExistsInternalHandler  = errors.New("BMCRevertNotExistsInternalHandler")
	ErrBMCRevertAlreadyExistsBMCPeriphery = errors.New("BMCRevertAlreadyExistsBMCPeriphery")
	ErrBMCRevertUnreachable               = errors.New("BMCRevertUnreachable:")
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
		ErrBTPRevertAlreadyExistRoute,
		ErrBTPRevertNotExistRoute,
		ErrBMCRevertInvalidSN,
		ErrBMCRevertParseFailure,
		ErrBMCRevertRxSeqLowerThanExpected,
		ErrBMCRevertRxSeqHigherThanExpected,
		ErrBMCRevertInvalidRxHeight,
		ErrBMCRevertNotExistsInternalHandler,
		ErrBMCRevertAlreadyExistsBMCPeriphery,
		ErrBMCRevertUnreachable,
	} {
		if strings.HasPrefix(msg, err.Error()) {
			return err
		}
	}
	return nil
}
