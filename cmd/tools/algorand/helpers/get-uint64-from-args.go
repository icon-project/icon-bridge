package helpers

import (
	"log"
	"os"
	"strconv"
)

func GetUint64FromArgs(index uint8, label string) uint64 {
	res, err := strconv.ParseUint(os.Args[index], 10, 64)

	if err != nil {
		log.Fatalf("Invalid argument for index %v, should be valid %s: %s\n", index, label, err)
	}

	return res
}