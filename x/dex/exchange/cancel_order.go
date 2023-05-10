package exchange

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/x/dex/keeper"
	"github.com/sei-protocol/sei-chain/x/dex/types"
	otrace "go.opentelemetry.io/otel/trace"
)

func CancelOrders(
	ctx sdk.Context, keeper *keeper.Keeper, contract types.ContractAddress, pair types.Pair,
	cancels []*types.Cancellation, tracer *otrace.Tracer,
	tracerCtx context.Context,
) {
	tracerCtx, span := (*tracer).Start(tracerCtx, "CancelOrders")
	defer span.End()
	span.AddEvent("CancelOrders")
	ctx.Logger().Info(fmt.Sprintf("CONTRACT-DDOS-2-DEBUG CancelOrders contract %s, pair %s, Num Cancellations %d\n", contract, pair, len(cancels)))
	for _, cancel := range cancels {
		cancelOrder(ctx, keeper, cancel, contract, pair, tracer, tracerCtx)
	}
}

func cancelOrder(ctx sdk.Context, keeper *keeper.Keeper, cancellation *types.Cancellation, contract types.ContractAddress, pair types.Pair, tracer *otrace.Tracer, tracerCtx context.Context) {
	tracerCtx, span := (*tracer).Start(tracerCtx, "cancelOrder")
	defer span.End()
	span.AddEvent("cancelOrder")
	getter, setter, deleter := keeper.GetLongOrderBookEntryByPrice, keeper.SetLongOrderBookEntry, keeper.RemoveLongBookByPrice
	if cancellation.PositionDirection == types.PositionDirection_SHORT {
		getter, setter, deleter = keeper.GetShortOrderBookEntryByPrice, keeper.SetShortOrderBookEntry, keeper.RemoveShortBookByPrice
	}
	entry, found := getter(ctx, string(contract), cancellation.Price, pair.PriceDenom, pair.AssetDenom)
	span.AddEvent("finished cancelOrder getter")
	if !found {
		return
	}
	newEntry := *entry.GetOrderEntry()
	newAllocations := []*types.Allocation{}
	newQuantity := sdk.ZeroDec()
	for _, allocation := range newEntry.Allocations {
		if allocation.OrderId != cancellation.Id {
			newAllocations = append(newAllocations, allocation)
			newQuantity = newQuantity.Add(allocation.Quantity)
		}
	}
	if newQuantity.IsZero() {
		deleter(ctx, string(contract), entry.GetPrice(), pair.PriceDenom, pair.AssetDenom)
		return
	}
	newEntry.Quantity = newQuantity
	newEntry.Allocations = newAllocations
	entry.SetEntry(&newEntry)
	setter(ctx, string(contract), entry)
	span.AddEvent("finished cancelOrder getsetterter")
}
