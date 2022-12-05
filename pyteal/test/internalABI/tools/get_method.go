package internalABI

import (
	"log"

	"github.com/algorand/go-algorand-sdk/abi"
)

func GetMethod(c *abi.Contract, name string) abi.Method {
	m, err := c.GetMethodByName(name)
	if err != nil {
		log.Fatalf("No method named: %s", name)
	}
	return m
}
