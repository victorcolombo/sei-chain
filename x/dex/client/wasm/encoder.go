package wasm

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/x/dex/types"
)

func EncodeDexPlaceOrders(rawMsg json.RawMessage) ([]sdk.Msg, error) {
	placeOrdersMsg := types.MsgPlaceOrders{}

	if true {
		return []sdk.Msg{}, sdkerrors.New("in custom endoer after unmarshal", 1, "xxx before unmarshall in placeorder")
	}

	if err := json.Unmarshal(rawMsg, &placeOrdersMsg); err != nil {
		return []sdk.Msg{}, err
	}

	if true {
		return []sdk.Msg{}, sdkerrors.New("in custom endoer after unmarshal", 1, "xxx after unmarshall in placeorder")
	}

	return []sdk.Msg{&placeOrdersMsg}, nil
}

func EncodeDexCancelOrders(rawMsg json.RawMessage) ([]sdk.Msg, error) {
	cancelOrdersMsg := types.MsgCancelOrders{}
	if err := json.Unmarshal(rawMsg, &cancelOrdersMsg); err != nil {
		return []sdk.Msg{}, err
	}
	return []sdk.Msg{&cancelOrdersMsg}, nil
}
