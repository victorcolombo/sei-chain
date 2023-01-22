package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// BroadcastTx attempts to generate, sign and broadcast a transaction with the
// given set of messages. It will also simulate gas requirements if necessary.
// It will return an error upon failure.
//
// Note, BroadcastTx is copied from the SDK except it removes a few unnecessary
// things like prompting for confirmation and printing the response. Instead,
// we return the TxResponse.
func BroadcastTx(clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	startTime := time.Now().UnixMicro()
	txf, err := prepareFactory(clientCtx, txf)
	if err != nil {
		return nil, err
	}

	prepareComplete := time.Now().UnixMicro()
	fmt.Printf("[Oracle BroadcastTx] Prepare latency is: %d\n", prepareComplete-startTime)

	// _, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
	// if err != nil {
	// 	return nil, err
	// }

	// txf = txf.WithGas(adjusted)
	txf = txf.WithGas(0)
	// txf = txf.WithFees("0usei")

	unsignedTx, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	// unsignedTx.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("usei", sdk.NewInt(2000))))

	if err = tx.Sign(txf, clientCtx.GetFromName(), unsignedTx, true); err != nil {
		return nil, err
	}

	signCompleted := time.Now().UnixMicro()

	fmt.Printf("[Oracle BroadcastTx] Sign latency is: %d\n", signCompleted-prepareComplete)

	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		return nil, err
	}

	result, err := clientCtx.BroadcastTx(txBytes)

	broadcastComplete := time.Now().UnixMicro()

	fmt.Printf("[Oracle BroadcastTx] Broadcast request latency is: %d\n", broadcastComplete-signCompleted)

	return result, err
}

// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
