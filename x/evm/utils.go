package evm

import (
	"errors"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/sei-protocol/sei-chain/wasmbinding/bindings"
)

func ToTransactionArgs(tx bindings.EthereumTx) *evmtypes.TransactionArgs {
	// create evmtypes.TransactionArgs

	// set default gas price/upper bound/lower bound

	// set chain id and nonce (sequence number) from auth module

	// encode bindings.EthereumTx to TransactionArgs.data ([]byte)
	abi, err := abi.JSON(strings.NewReader(tx.MethodJson))
	if err != nil {
		//
	}

	// example below, we need to go over the reflect types and do type assertion
	var j interface{} = "1000"
	for i := 0; i < len(abi.Methods["send"].Inputs); i++ {
		switch abi.Methods["send"].Inputs[i].Type.GetType() {
			case reflect.TypeOf(&big.Int{}):
				n := new(big.Int) 
				j, _ = n.SetString(j.(string), 0)
		}
	}
	if _, err := abi.Pack("send", j); err != nil {
		//
	}

	return nil
}

// REFERENCE below, need to remove later
/*
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *ethtypes.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big         `json:"chainId,omitempty"`
}

// SetTxDefaults populates tx message with default values in case they are not
// provided on the args
func SetTxDefaults(args TransactionArgs) (evmtypes.TransactionArgs, error) {
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return args, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}

	head := b.CurrentHeader()
	if head == nil {
		return args, errors.New("latest header is nil")
	}

	// If user specifies both maxPriorityfee and maxFee, then we do not
	// need to consult the chain for defaults. It's definitely a London tx.
	if args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil {
		// In this clause, user left some fields unspecified.
		if head.BaseFee != nil && args.GasPrice == nil {
			if args.MaxPriorityFeePerGas == nil {
				tip, err := b.SuggestGasTipCap(head.BaseFee)
				if err != nil {
					return args, err
				}
				args.MaxPriorityFeePerGas = (*hexutil.Big)(tip)
			}

			if args.MaxFeePerGas == nil {
				gasFeeCap := new(big.Int).Add(
					(*big.Int)(args.MaxPriorityFeePerGas),
					new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
				)
				args.MaxFeePerGas = (*hexutil.Big)(gasFeeCap)
			}

			if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
				return args, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
			}

		} else {
			if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
				return args, errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
			}

			if args.GasPrice == nil {
				price, err := b.SuggestGasTipCap(head.BaseFee)
				if err != nil {
					return args, err
				}
				if head.BaseFee != nil {
					// The legacy tx gas price suggestion should not add 2x base fee
					// because all fees are consumed, so it would result in a spiral
					// upwards.
					price.Add(price, head.BaseFee)
				}
				args.GasPrice = (*hexutil.Big)(price)
			}
		}
	} else {
		// Both maxPriorityfee and maxFee set by caller. Sanity-check their internal relation
		if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return args, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
		}
	}

	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		// get the nonce from the account retriever
		// ignore error in case tge account doesn't exist yet
		nonce, _ := b.getAccountNonce(*args.From, true, 0, b.logger)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}

	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return args, errors.New("both 'data' and 'input' are set and not equal. Please use 'input' to pass transaction call data")
	}

	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}

		if len(input) == 0 {
			return args, errors.New("contract creation without any data provided")
		}
	}

	if args.Gas == nil {
		// For backwards-compatibility reason, we try both input and data
		// but input is preferred.
		input := args.Input
		if input == nil {
			input = args.Data
		}

		callArgs := evmtypes.TransactionArgs{
			From:                 args.From,
			To:                   args.To,
			Gas:                  args.Gas,
			GasPrice:             args.GasPrice,
			MaxFeePerGas:         args.MaxFeePerGas,
			MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
			Value:                args.Value,
			Data:                 input,
			AccessList:           args.AccessList,
		}

		blockNr := rpctypes.NewBlockNumber(big.NewInt(0))
		estimated, err := b.EstimateGas(callArgs, &blockNr)
		if err != nil {
			return args, err
		}
		args.Gas = &estimated
		b.logger.Debug("estimate gas usage automatically", "gas", args.Gas)
	}

	if args.ChainID == nil {
		args.ChainID = (*hexutil.Big)(b.chainID)
	}

	return args, nil
}

// ToTransaction converts the arguments to an ethereum transaction.
// This assumes that setTxDefaults has been called.
func (args *TransactionArgs) ToTransaction() *evmtypes.MsgEthereumTx {
	var (
		chainID, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas sdkmath.Int
		gas, nonce                                                   uint64
		from, to                                                     string
	)

	// Set sender address or use zero address if none specified.
	if args.ChainID != nil {
		chainID = sdkmath.NewIntFromBigInt(args.ChainID.ToInt())
	}

	if args.Nonce != nil {
		nonce = uint64(*args.Nonce)
	}

	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}

	if args.GasPrice != nil {
		gasPrice = sdkmath.NewIntFromBigInt(args.GasPrice.ToInt())
	}

	if args.MaxFeePerGas != nil {
		maxFeePerGas = sdkmath.NewIntFromBigInt(args.MaxFeePerGas.ToInt())
	}

	if args.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = sdkmath.NewIntFromBigInt(args.MaxPriorityFeePerGas.ToInt())
	}

	if args.Value != nil {
		value = sdkmath.NewIntFromBigInt(args.Value.ToInt())
	}

	if args.To != nil {
		to = args.To.Hex()
	}

	var data TxData
	switch {
	case args.MaxFeePerGas != nil:
		al := AccessList{}
		if args.AccessList != nil {
			al = NewAccessList(args.AccessList)
		}

		data = &DynamicFeeTx{
			To:        to,
			ChainID:   &chainID,
			Nonce:     nonce,
			GasLimit:  gas,
			GasFeeCap: &maxFeePerGas,
			GasTipCap: &maxPriorityFeePerGas,
			Amount:    &value,
			Data:      args.GetData(),
			Accesses:  al,
		}
	case args.AccessList != nil:
		data = &AccessListTx{
			To:       to,
			ChainID:  &chainID,
			Nonce:    nonce,
			GasLimit: gas,
			GasPrice: &gasPrice,
			Amount:   &value,
			Data:     args.GetData(),
			Accesses: NewAccessList(args.AccessList),
		}
	default:
		data = &LegacyTx{
			To:       to,
			Nonce:    nonce,
			GasLimit: gas,
			GasPrice: &gasPrice,
			Amount:   &value,
			Data:     args.GetData(),
		}
	}

	any, err := PackTxData(data)
	if err != nil {
		return nil
	}

	if args.From != nil {
		from = args.From.Hex()
	}

	msg := MsgEthereumTx{
		Data: any,
		From: from,
	}
	msg.Hash = msg.AsTransaction().Hash().Hex()
	return &msg
}
*/