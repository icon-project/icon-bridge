package algo

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

func deployContract(ctx context.Context, algodAccess []string, tealPath [2]string, account crypto.Account) (uint64, error) {
	client, err := algod.MakeClient(algodAccess[0], algodAccess[1])
	if err != nil {
		return 0, fmt.Errorf("Bmc couldn't create algod: %w", err)
	}
	params, err := client.SuggestedParams().Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("Error getting params: %w", err)
	}

	// call python script to compile bmc
	_, err = exec.Command(bmcCompilePyPath).Output()
	if err != nil {
		return 0, fmt.Errorf("Failed to compile bmc.py: %w", err)
	}

	approvalFile, err := filepath.Abs(contractDir + tealPath[0])
	if err != nil {
		return 0, fmt.Errorf("Couldn't find approval file: %w", err)
	}
	approvalProgram, err := compileTeal(ctx, client, approvalFile)
	if err != nil {
		return 0, fmt.Errorf("Approval compile err: %w", err)
	}

	clearFile, err := filepath.Abs(contractDir + tealPath[1])
	if err != nil {
		return 0, fmt.Errorf("Couldn't find clear file: %w", err)
	}
	clearProgram, err := compileTeal(ctx, client, clearFile)
	if err != nil {
		return 0, fmt.Errorf("Clear compile err: %w", err)
	}

	txn, err := future.MakeApplicationCreateTx(
		false,
		approvalProgram,
		clearProgram,
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
		[][]byte{},
		nil,
		nil,
		nil,
		params,
		account.Address,
		nil,
		types.Digest{},
		[32]byte{},
		types.Address{},
	)
	if err != nil {
		return 0, fmt.Errorf("Failed to make bmc: %w", err)
	}

	txID, signedTxn, err := crypto.SignTransaction(account.PrivateKey, txn)
	if err != nil {
		return 0, fmt.Errorf("Failed to sign transaction: %w", err)
	}
	_, err = client.SendRawTransaction(signedTxn).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed to send transaction: %w", err)
	}
	deployRes, err := future.WaitForConfirmation(client, txID, waitRounds, ctx)
	if err != nil {
		return 0, fmt.Errorf("Error waiting for confirmation: %w", err)
	}
	return deployRes.ApplicationIndex, nil
}

func compileTeal(ctx context.Context, client *algod.Client, filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}
	compileResponse, err := client.TealCompile(content).Do(ctx)
	if err != nil {
		return []byte{}, err
	}

	decodedProgram, err := base64.StdEncoding.DecodeString(compileResponse.Result)
	if err != nil {
		return []byte{}, err
	}
	return decodedProgram, nil
}
