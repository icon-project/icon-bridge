package lifecheck_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func Test_Lifecheck(t *testing.T) {
    algodAddress := os.Getenv("ALGOD_ADDRESS")

    if (algodAddress == "") {
        algodAddress = "http://localhost:4001"
    }

    client, err := algod.MakeClient(algodAddress, algodToken)
    
    if err != nil {
        os.Exit(1)
    }

    for i := 0 ; i < 10; i++{
        status, _ := client.Status().Do(context.Background())

        if (i >= 5 && status.LastRound == 0) {
            t.Fatal("Can't get last round")
        }

        fmt.Println(status.LastRound, " -- " , time.Now())
        time.Sleep(time.Duration(500)*time.Millisecond)
    }
}