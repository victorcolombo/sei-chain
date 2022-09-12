package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/sei-protocol/sei-chain/wasmbinding/bindings"
	"github.com/sei-protocol/sei-chain/x/evm"
)

func EncodeEthereumMsg(ctx sdk.Context, accountKeeper *authkeeper.AccountKeeper, rawMsg json.RawMessage, sender sdk.AccAddress) ([]sdk.Msg, error) {
	encodedEthereumTxMsg := bindings.EthereumTx{}
	if err := json.Unmarshal(rawMsg, &encodedEthereumTxMsg); err != nil {
		return []sdk.Msg{}, err
	}

	ethereumTransactionArgs := evm.ToTransactionArgs(ctx, accountKeeper, encodedEthereumTxMsg)
	msgEthereumTx := ethereumTransactionArgs.ToTransaction()
	return []sdk.Msg{msgEthereumTx}, nil
}
