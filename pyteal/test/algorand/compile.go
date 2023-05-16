package algorand

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func Compile(client *algod.Client, teal []byte) (program []byte, err error) {
	var response models.CompileResponse
	response, err = client.TealCompile(teal).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("compilation failed: %+v\n ", err)
		return
	}

	program, err = base64.StdEncoding.DecodeString(response.Result)
	if err != nil {
		err = fmt.Errorf("failed to base64 decode compiled program: %s", err)
		return
	}
	return
}
