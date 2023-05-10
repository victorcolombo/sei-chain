package contract

import (
	"context"
	"fmt"
	"sync"
	"time"

	otrace "go.opentelemetry.io/otel/trace"

	"github.com/cosmos/cosmos-sdk/telemetry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/store/whitelist/multi"
	"github.com/sei-protocol/sei-chain/utils/datastructures"
	"github.com/sei-protocol/sei-chain/x/dex/exchange"
	"github.com/sei-protocol/sei-chain/x/dex/keeper"
	dexkeeperabci "github.com/sei-protocol/sei-chain/x/dex/keeper/abci"
	dexkeeperutils "github.com/sei-protocol/sei-chain/x/dex/keeper/utils"
	"github.com/sei-protocol/sei-chain/x/dex/types"
	dexutils "github.com/sei-protocol/sei-chain/x/dex/utils"
	"go.opentelemetry.io/otel/attribute"
)

func CallPreExecutionHooks(
	ctx context.Context,
	sdkCtx sdk.Context,
	contractAddr string,
	dexkeeper *keeper.Keeper,
	registeredPairs []types.Pair,
	tracer *otrace.Tracer,
) error {
	spanCtx, span := (*tracer).Start(ctx, "PreExecutionHooks")
	defer span.End()
	span.SetAttributes(attribute.String("contract", contractAddr))
	abciWrapper := dexkeeperabci.KeeperWrapper{Keeper: dexkeeper}
	if err := abciWrapper.HandleEBCancelOrders(spanCtx, sdkCtx, tracer, contractAddr, registeredPairs); err != nil {
		return err
	}
	if err := abciWrapper.HandleEBPlaceOrders(spanCtx, sdkCtx, tracer, contractAddr, registeredPairs); err != nil {
		return err
	}
	return nil
}

func ExecutePair(
	ctx sdk.Context,
	contractAddr string,
	pair types.Pair,
	dexkeeper *keeper.Keeper,
	orderbook *types.OrderBook,
	tracer *otrace.Tracer,
	tracerCtx context.Context,
) []*types.SettlementEntry {
	tracerCtx, span := (*tracer).Start(tracerCtx, "ExecutePair")
	defer span.End()
	typedContractAddr := types.ContractAddress(contractAddr)
	typedPairStr := types.GetPairString(&pair)

	// First cancel orders
	span.AddEvent("cancelForPair")
	cancelForPair(ctx, dexkeeper, typedContractAddr, pair, tracer, tracerCtx)
	// Add all limit orders to the orderbook
	span.AddEvent("getMemState")
	orders := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr)
	span.AddEvent("GetLimitOrders_LONG")
	limitBuys := orders.GetLimitOrders(types.PositionDirection_LONG)
	span.AddEvent("GetLimitOrders_SHORT")
	limitSells := orders.GetLimitOrders(types.PositionDirection_SHORT)
	span.AddEvent("exchange.AddOutstandingLimitOrdersToOrderbook")
	exchange.AddOutstandingLimitOrdersToOrderbook(ctx, dexkeeper, limitBuys, limitSells)
	// Fill market orders
	span.AddEvent("matchMarketOrderForPair")
	marketOrderOutcome := matchMarketOrderForPair(ctx, typedContractAddr, typedPairStr, orderbook, tracer, tracerCtx)
	// Fill limit orders
	span.AddEvent("MatchLimitOrders")
	limitOrderOutcome := exchange.MatchLimitOrders(ctx, orderbook)
	span.AddEvent("marketOrderOutcome.Merge")
	totalOutcome := marketOrderOutcome.Merge(&limitOrderOutcome)
	span.AddEvent("UpdateTriggeredOrderForPair.Merge")
	UpdateTriggeredOrderForPair(ctx, typedContractAddr, typedPairStr, dexkeeper, totalOutcome)
	span.AddEvent("SetPriceStateFromExecutionOutcome")
	dexkeeperutils.SetPriceStateFromExecutionOutcome(ctx, dexkeeper, typedContractAddr, pair, totalOutcome)
	span.AddEvent("done")
	return totalOutcome.Settlements
}

func cancelForPair(
	ctx sdk.Context,
	keeper *keeper.Keeper,
	contractAddress types.ContractAddress,
	pair types.Pair,
	tracer *otrace.Tracer,
	tracerCtx context.Context,
) {
	tracerCtx, span := (*tracer).Start(tracerCtx, "cancelForPair")
	defer span.End()
	span.AddEvent("cancelForPair")
	cancels := dexutils.GetMemState(ctx.Context()).GetBlockCancels(ctx, contractAddress, types.GetPairString(&pair))
	span.AddEvent("GetMemState")
	exchange.CancelOrders(ctx, keeper, contractAddress, pair, cancels.Get(), tracer, tracerCtx)
	span.AddEvent("CancelOrders")
}

func matchMarketOrderForPair(
	ctx sdk.Context,
	typedContractAddr types.ContractAddress,
	typedPairStr types.PairString,
	orderbook *types.OrderBook,
	tracer *otrace.Tracer,
	tracerCtx context.Context,
) exchange.ExecutionOutcome {
	tracerCtx, span := (*tracer).Start(tracerCtx, "matchMarketOrderForPair")
	defer span.End()
	span.AddEvent("matchMarketOrderForPair")
	orders := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr)
	marketBuys := orders.GetSortedMarketOrders(types.PositionDirection_LONG)
	marketSells := orders.GetSortedMarketOrders(types.PositionDirection_SHORT)
	ctx.Logger().Info(fmt.Sprintf("CONTRACT-DDOS-2-DEBUG matchMarketOrderForPair, contract %s, pair %s, num marketBuys %d\n", typedContractAddr, typedPairStr, len(marketBuys)))
	ctx.Logger().Info(fmt.Sprintf("CONTRACT-DDOS-2-DEBUG matchMarketOrderForPair, contract %s, pair %s, num marketSells %d\n", typedContractAddr, typedPairStr, len(marketSells)))
	span.AddEvent("finished get sorted market buys, market sells")
	marketBuyOutcome := exchange.MatchMarketOrders(
		ctx,
		marketBuys,
		orderbook.Shorts,
		types.PositionDirection_LONG,
		orders,
	)
	marketSellOutcome := exchange.MatchMarketOrders(
		ctx,
		marketSells,
		orderbook.Longs,
		types.PositionDirection_SHORT,
		orders,
	)
	span.AddEvent("finished MatchMarket Orders")
	mergedMarketBuyOutcome := marketBuyOutcome.Merge(&marketSellOutcome)
	span.AddEvent("finished mergedMarketBuyOutcome")
	return mergedMarketBuyOutcome
}

func MoveTriggeredOrderForPair(
	ctx sdk.Context,
	typedContractAddr types.ContractAddress,
	typedPairStr types.PairString,
	dexkeeper *keeper.Keeper,
) {
	priceDenom, assetDenom := types.GetPriceAssetString(typedPairStr)
	triggeredOrders := dexkeeper.GetAllTriggeredOrdersForPair(ctx, string(typedContractAddr), priceDenom, assetDenom)
	for i, order := range triggeredOrders {
		if order.TriggerStatus {
			if order.OrderType == types.OrderType_STOPLOSS {
				triggeredOrders[i].OrderType = types.OrderType_MARKET
			} else if order.OrderType == types.OrderType_STOPLIMIT {
				triggeredOrders[i].OrderType = types.OrderType_LIMIT
			}
			dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr).Add(&triggeredOrders[i])
			dexkeeper.RemoveTriggeredOrder(ctx, string(typedContractAddr), order.Id, priceDenom, assetDenom)
		}
	}
}

func UpdateTriggeredOrderForPair(
	ctx sdk.Context,
	typedContractAddr types.ContractAddress,
	typedPairStr types.PairString,
	dexkeeper *keeper.Keeper,
	totalOutcome exchange.ExecutionOutcome,
) {
	// update existing trigger orders
	priceDenom, assetDenom := types.GetPriceAssetString(typedPairStr)
	triggeredOrders := dexkeeper.GetAllTriggeredOrdersForPair(ctx, string(typedContractAddr), priceDenom, assetDenom)
	for i, order := range triggeredOrders {
		if order.PositionDirection == types.PositionDirection_LONG && order.TriggerPrice.LTE(totalOutcome.MaxPrice) {
			triggeredOrders[i].TriggerStatus = true
			dexkeeper.SetTriggeredOrder(ctx, string(typedContractAddr), triggeredOrders[i], priceDenom, assetDenom)
		} else if order.PositionDirection == types.PositionDirection_SHORT && order.TriggerPrice.GTE(totalOutcome.MinPrice) {
			triggeredOrders[i].TriggerStatus = true
			dexkeeper.SetTriggeredOrder(ctx, string(typedContractAddr), triggeredOrders[i], priceDenom, assetDenom)
		}
	}

	// update triggered orders in cache
	orders := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr)
	cacheTriggeredOrders := orders.GetTriggeredOrders()
	for i, order := range cacheTriggeredOrders {
		if order.PositionDirection == types.PositionDirection_LONG && order.TriggerPrice.LTE(totalOutcome.MaxPrice) {
			cacheTriggeredOrders[i].TriggerStatus = true
		} else if order.PositionDirection == types.PositionDirection_SHORT && order.TriggerPrice.GTE(totalOutcome.MinPrice) {
			cacheTriggeredOrders[i].TriggerStatus = true
		}
		dexkeeper.SetTriggeredOrder(ctx, string(typedContractAddr), *cacheTriggeredOrders[i], priceDenom, assetDenom)
	}
}

func GetMatchResults(
	ctx sdk.Context,
	typedContractAddr types.ContractAddress,
	typedPairStr types.PairString,
) ([]*types.Order, []*types.Cancellation) {
	orderResults := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr).Get()
	cancelResults := dexutils.GetMemState(ctx.Context()).GetBlockCancels(ctx, typedContractAddr, typedPairStr).Get()
	return orderResults, cancelResults
}

func GetOrderIDToSettledQuantities(settlements []*types.SettlementEntry) map[uint64]sdk.Dec {
	res := map[uint64]sdk.Dec{}
	for _, settlement := range settlements {
		if _, ok := res[settlement.OrderId]; !ok {
			res[settlement.OrderId] = sdk.ZeroDec()
		}
		res[settlement.OrderId] = res[settlement.OrderId].Add(settlement.Quantity)
	}
	return res
}

func ExecutePairsInParallel(
	sdkCtx sdk.Context,
	contractAddr string,
	dexkeeper *keeper.Keeper,
	registeredPairs []types.Pair,
	orderBooks *datastructures.TypedSyncMap[types.PairString, *types.OrderBook],
	tracer *otrace.Tracer,
	ctx context.Context,
) []*types.SettlementEntry {
	tracerCtx, span := (*tracer).Start(ctx, "ExecutePairsInParallel")
	defer span.End()

	typedContractAddr := types.ContractAddress(contractAddr)
	orderResults := []*types.Order{}
	cancelResults := []*types.Cancellation{}
	settlements := []*types.SettlementEntry{}

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, pair := range registeredPairs {
		wg.Add(1)

		pair := pair
		pairCtx := sdkCtx.WithMultiStore(multi.NewStore(sdkCtx.MultiStore(), GetPerPairWhitelistMap(contractAddr, pair))).WithEventManager(sdk.NewEventManager())
		go func() {
			_, pairSpan := (*tracer).Start(tracerCtx, "execeute pair")
			defer pairSpan.End()

			defer wg.Done()
			pairCopy := pair
			pairStr := types.GetPairString(&pairCopy)
			pairSpan.AddEvent("MoveTriggeredOrderForPair")
			MoveTriggeredOrderForPair(sdkCtx, typedContractAddr, pairStr, dexkeeper)
			orderbook, found := orderBooks.Load(pairStr)
			if !found {
				panic(fmt.Sprintf("Orderbook not found for %s", pairStr))
			}
			pairSpan.AddEvent("ExecutePair")
			pairSettlements := ExecutePair(pairCtx, contractAddr, pair, dexkeeper, orderbook, tracer, tracerCtx)
			pairSpan.AddEvent("GetOrderIDToSettledQuantities")
			orderIDToSettledQuantities := GetOrderIDToSettledQuantities(pairSettlements)
			pairSpan.AddEvent("PrepareCancelUnfulfilledMarketOrders")
			PrepareCancelUnfulfilledMarketOrders(pairCtx, typedContractAddr, pairStr, orderIDToSettledQuantities)

			pairSpan.AddEvent("Waiting for lock")
			mu.Lock()
			defer mu.Unlock()
			pairSpan.AddEvent("Get Match Results")
			orders, cancels := GetMatchResults(sdkCtx, typedContractAddr, types.GetPairString(&pairCopy))
			pairSpan.AddEvent("Finish")
			defer pairSpan.AddEvent("Finish")
			orderResults = append(orderResults, orders...)
			cancelResults = append(cancelResults, cancels...)
			settlements = append(settlements, pairSettlements...)
			// ordering of events doesn't matter since events aren't part of consensus
			sdkCtx.EventManager().EmitEvents(pairCtx.EventManager().Events())
		}()
	}
	wg.Wait()
	dexkeeper.SetMatchResult(sdkCtx, contractAddr, types.NewMatchResult(orderResults, cancelResults, settlements))

	return settlements
}

func HandleExecutionForContract(
	ctx context.Context,
	sdkCtx sdk.Context,
	contract types.ContractInfoV2,
	dexkeeper *keeper.Keeper,
	registeredPairs []types.Pair,
	orderBooks *datastructures.TypedSyncMap[types.PairString, *types.OrderBook],
	tracer *otrace.Tracer,
) ([]*types.SettlementEntry, error) {
	tracerCtx, span := (*tracer).Start(ctx, "HandleExecutionForContract")
	defer span.End()
	executionStart := time.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, executionStart, "handle_execution_for_contract_ms")
	contractAddr := contract.ContractAddr

	// Call contract hooks so that contracts can do internal bookkeeping
	span.AddEvent("HandleExecutionForContract: CallPreExecutionHooks start")
	if err := CallPreExecutionHooks(ctx, sdkCtx, contractAddr, dexkeeper, registeredPairs, tracer); err != nil {
		return []*types.SettlementEntry{}, err
	}
	span.AddEvent("HandleExecutionForContract: CallPreExecutionHooks end")

	settlements := ExecutePairsInParallel(
		sdkCtx,
		contractAddr,
		dexkeeper,
		registeredPairs,
		orderBooks,
		tracer,
		tracerCtx,
	)
	span.AddEvent("HandleExecutionForContract: ExecutePairsInParallel end")
	defer EmitSettlementMetrics(settlements)

	return settlements, nil
}

// Emit metrics for settlements
func EmitSettlementMetrics(settlements []*types.SettlementEntry) {
	if len(settlements) > 0 {
		telemetry.ModuleSetGauge(
			types.ModuleName,
			float32(len(settlements)),
			"num_settlements",
		)
		var totalQuantity int
		for _, s := range settlements {
			totalQuantity += s.Quantity.Size()
			telemetry.IncrCounter(
				1,
				"num_settlements_order_type_"+s.GetOrderType(),
			)
			telemetry.IncrCounter(
				1,
				"num_settlements_position_direction"+s.GetPositionDirection(),
			)
			telemetry.IncrCounter(
				1,
				"num_settlements_asset_denom_"+s.GetAssetDenom(),
			)
			telemetry.IncrCounter(
				1,
				"num_settlements_price_denom_"+s.GetPriceDenom(),
			)
		}
		telemetry.ModuleSetGauge(
			types.ModuleName,
			float32(totalQuantity),
			"num_total_order_quantity_in_settlements",
		)
	}
}
