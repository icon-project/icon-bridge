package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetFileVar(cacheDir string, filename string) string {
	file, err := os.Open(cacheDir + filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	// read file contents as byte slice
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// convert byte slice to string
	return string(byteValue)
}