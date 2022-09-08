package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/sei-protocol/sei-chain/wasmbinding/bindings"
	"github.com/sei-protocol/sei-chain/x/tokenfactory/types"
	"github.com/sei-protocol/sei-chain/x/evm"
)

func EncodeEthereumMsg(rawMsg json.RawMessage, sender sdk.AccAddress) ([]sdk.Msg, error) {
	encodedEthereumTxMsg := bindings.EthereumTx{}
	if err := json.Unmarshal(rawMsg, &encodedEthereumTxMsg); err != nil {
		return []sdk.Msg{}, err
	}

	ethereumTransactionArgs := evm.ToTransactionArgs(encodedEthereumTxMsg)
	msgEthereumTx := ethereumTransactionArgs.ToTransaction()
	return []sdk.Msg{msgEthereumTx}, nil
}
