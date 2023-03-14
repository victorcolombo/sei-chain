package contract

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	dextypesutils "github.com/sei-protocol/sei-chain/x/dex/types/utils"
	dextypeswasm "github.com/sei-protocol/sei-chain/x/dex/types/wasm"
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

var TotalExecuteLatency = atomic.Int64{}
var TotalCancelLatency = atomic.Int64{}
var TotalGetLimitOrdersLatency = atomic.Int64{}
var TotalAddLimitOrderLatency = atomic.Int64{}
var TotalFillMarketOrderLatency = atomic.Int64{}
var TotalFillLimitOrderLatency = atomic.Int64{}
var TotalMergeLatency = atomic.Int64{}
var TotalUpdateOrderLatency = atomic.Int64{}

func ExecutePair(
	ctx sdk.Context,
	contractAddr string,
	pair types.Pair,
	dexkeeper *keeper.Keeper,
	orderbook *types.OrderBook,
) []*types.SettlementEntry {
	startTime := time.Now().UnixMicro()
	typedContractAddr := dextypesutils.ContractAddress(contractAddr)
	typedPairStr := dextypesutils.GetPairString(&pair)

	// First cancel orders
	cancelForPair(ctx, typedContractAddr, typedPairStr, orderbook)
	endCancelOrderTime := time.Now().UnixMicro()
	TotalCancelLatency.Add(endCancelOrderTime - startTime)
	// Add all limit orders to the orderbook
	orders := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr)
	limitBuys := orders.GetLimitOrders(types.PositionDirection_LONG)
	limitSells := orders.GetLimitOrders(types.PositionDirection_SHORT)
	endGetLimitOrdersTime := time.Now().UnixMicro()
	TotalGetLimitOrdersLatency.Add(endGetLimitOrdersTime - endCancelOrderTime)

	exchange.AddOutstandingLimitOrdersToOrderbook(orderbook, limitBuys, limitSells)
	endAddOrderTime := time.Now().UnixMicro()
	TotalAddLimitOrderLatency.Add(endAddOrderTime - endGetLimitOrdersTime)
	// Fill market orders
	marketOrderOutcome := matchMarketOrderForPair(ctx, typedContractAddr, typedPairStr, orderbook)
	endFillMarketTime := time.Now().UnixMicro()
	TotalFillMarketOrderLatency.Add(endFillMarketTime - endAddOrderTime)
	// Fill limit orders
	limitOrderOutcome := exchange.MatchLimitOrders(ctx, orderbook)
	endFillLimitTime := time.Now().UnixMicro()
	TotalFillLimitOrderLatency.Add(endFillLimitTime - endFillMarketTime)
	// Merge order outcome
	totalOutcome := marketOrderOutcome.Merge(&limitOrderOutcome)
	endMergeTime := time.Now().UnixMicro()
	TotalMergeLatency.Add(endMergeTime - endFillLimitTime)

	// Update
	UpdateTriggeredOrderForPair(ctx, typedContractAddr, typedPairStr, dexkeeper, totalOutcome)
	dexkeeperutils.SetPriceStateFromExecutionOutcome(ctx, dexkeeper, typedContractAddr, pair, totalOutcome)
	dexkeeperutils.FlushOrderbook(ctx, dexkeeper, typedContractAddr, orderbook)
	endTime := time.Now().UnixMicro()
	TotalUpdateOrderLatency.Add(endTime - endMergeTime)
	TotalExecuteLatency.Add(endTime - startTime)
	return totalOutcome.Settlements
}

func cancelForPair(
	ctx sdk.Context,
	typedContractAddr dextypesutils.ContractAddress,
	typedPairStr dextypesutils.PairString,
	orderbook *types.OrderBook,
) {
	cancels := dexutils.GetMemState(ctx.Context()).GetBlockCancels(ctx, typedContractAddr, typedPairStr)
	exchange.CancelOrders(cancels.Get(), orderbook)
}

var TotalGetBlockOrderLatency = atomic.Int64{}
var TotalGetLongOrdersLatency = atomic.Int64{}
var TotalGetShortOrdersLatency = atomic.Int64{}
var TotalMarketBuyLatency = atomic.Int64{}
var TotalMarketSellLatency = atomic.Int64{}

func matchMarketOrderForPair(
	ctx sdk.Context,
	typedContractAddr dextypesutils.ContractAddress,
	typedPairStr dextypesutils.PairString,
	orderbook *types.OrderBook,
) exchange.ExecutionOutcome {
	startTime := time.Now().UnixMicro()
	orders := dexutils.GetMemState(ctx.Context()).GetBlockOrders(ctx, typedContractAddr, typedPairStr)
	endGetBlockOrderTime := time.Now().UnixMicro()
	TotalGetBlockOrderLatency.Add(endGetBlockOrderTime - startTime)
	marketBuys := orders.GetSortedMarketOrders(types.PositionDirection_LONG, true)
	endGetLongOrdersTime := time.Now().UnixMicro()
	TotalGetLongOrdersLatency.Add(endGetLongOrdersTime - endGetBlockOrderTime)
	marketSells := orders.GetSortedMarketOrders(types.PositionDirection_SHORT, true)
	endGetShortOrdersTime := time.Now().UnixMicro()
	TotalGetShortOrdersLatency.Add(endGetShortOrdersTime - endGetLongOrdersTime)
	marketBuyOutcome := exchange.MatchMarketOrders(
		ctx,
		marketBuys,
		orderbook.Shorts,
		types.PositionDirection_LONG,
		orders,
	)
	endMarketBuyTime := time.Now().UnixMicro()
	TotalMarketBuyLatency.Add(endMarketBuyTime - endGetShortOrdersTime)
	marketSellOutcome := exchange.MatchMarketOrders(
		ctx,
		marketSells,
		orderbook.Longs,
		types.PositionDirection_SHORT,
		orders,
	)
	endMarketSellTime := time.Now().UnixMicro()
	TotalMarketSellLatency.Add(endMarketSellTime - endMarketBuyTime)
	return marketBuyOutcome.Merge(&marketSellOutcome)
}

func MoveTriggeredOrderForPair(
	ctx sdk.Context,
	typedContractAddr dextypesutils.ContractAddress,
	typedPairStr dextypesutils.PairString,
	dexkeeper *keeper.Keeper,
) {
	priceDenom, assetDenom := dextypesutils.GetPriceAssetString(typedPairStr)
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
	typedContractAddr dextypesutils.ContractAddress,
	typedPairStr dextypesutils.PairString,
	dexkeeper *keeper.Keeper,
	totalOutcome exchange.ExecutionOutcome,
) {
	// update existing trigger orders
	priceDenom, assetDenom := dextypesutils.GetPriceAssetString(typedPairStr)
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
	typedContractAddr dextypesutils.ContractAddress,
	typedPairStr dextypesutils.PairString,
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

func ExecutePairsInParallel(ctx sdk.Context, contractAddr string, dexkeeper *keeper.Keeper, registeredPairs []types.Pair, orderBooks *datastructures.TypedSyncMap[dextypesutils.PairString, *types.OrderBook]) ([]*types.SettlementEntry, []*types.Cancellation) {

	typedContractAddr := dextypesutils.ContractAddress(contractAddr)
	orderResults := []*types.Order{}
	cancelResults := []*types.Cancellation{}
	settlements := []*types.SettlementEntry{}

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, pair := range registeredPairs {
		wg.Add(1)

		pair := pair
		pairCtx := ctx.WithMultiStore(multi.NewStore(ctx.MultiStore(), GetPerPairWhitelistMap(contractAddr, pair))).WithEventManager(sdk.NewEventManager())
		go func() {
			//startTime := time.Now().UnixMicro()
			defer wg.Done()
			pairCopy := pair
			pairStr := dextypesutils.GetPairString(&pairCopy)
			MoveTriggeredOrderForPair(ctx, typedContractAddr, pairStr, dexkeeper)
			orderbook, found := orderBooks.Load(pairStr)
			//endLoadingOrderBookTime := time.Now().UnixMicro()

			if !found {
				panic(fmt.Sprintf("Orderbook not found for %s", pairStr))
			}

			pairSettlements := ExecutePair(pairCtx, contractAddr, pair, dexkeeper, orderbook.DeepCopy())
			//endExecutingPairTime := time.Now().UnixMicro()

			orderIDToSettledQuantities := GetOrderIDToSettledQuantities(pairSettlements)
			PrepareCancelUnfulfilledMarketOrders(pairCtx, typedContractAddr, pairStr, orderIDToSettledQuantities)
			//endPrepareCancelTime := time.Now().UnixMicro()

			mu.Lock()
			defer mu.Unlock()
			orders, cancels := GetMatchResults(ctx, typedContractAddr, dextypesutils.GetPairString(&pairCopy))
			orderResults = append(orderResults, orders...)
			cancelResults = append(cancelResults, cancels...)
			settlements = append(settlements, pairSettlements...)
			// ordering of events doesn't matter since events aren't part of consensus
			ctx.EventManager().EmitEvents(pairCtx.EventManager().Events())
			//endTime := time.Now().UnixMicro()
			//ctx.Logger().Info(fmt.Sprintf("[SeiChain-Debug] Execution pair %s total time %d, loadOrderbook latency: %d, ExecutePair latency: %d, PrepareCancel latency: %d, GetMatchResults latency: %d", pairStr, endTime-startTime, endLoadingOrderBookTime-startTime, endExecutingPairTime-endLoadingOrderBookTime, endPrepareCancelTime-endExecutingPairTime, endTime-endPrepareCancelTime))
		}()
	}
	wg.Wait()
	dexkeeper.SetMatchResult(ctx, contractAddr, types.NewMatchResult(orderResults, cancelResults, settlements))
	return settlements, cancelResults
}

func HandleExecutionForContract(
	ctx context.Context,
	sdkCtx sdk.Context,
	contract types.ContractInfoV2,
	dexkeeper *keeper.Keeper,
	registeredPairs []types.Pair,
	orderBooks *datastructures.TypedSyncMap[dextypesutils.PairString, *types.OrderBook],
	tracer *otrace.Tracer,
) (map[string]dextypeswasm.ContractOrderResult, []*types.SettlementEntry, error) {
	executionStart := time.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, executionStart, "handle_execution_for_contract_ms")
	contractAddr := contract.ContractAddr
	typedContractAddr := dextypesutils.ContractAddress(contractAddr)
	orderResults := map[string]dextypeswasm.ContractOrderResult{}

	// Call contract hooks so that contracts can do internal bookkeeping
	if err := CallPreExecutionHooks(ctx, sdkCtx, contractAddr, dexkeeper, registeredPairs, tracer); err != nil {
		return orderResults, []*types.SettlementEntry{}, err
	}
	settlements, cancellations := ExecutePairsInParallel(sdkCtx, contractAddr, dexkeeper, registeredPairs, orderBooks)
	defer EmitSettlementMetrics(settlements)
	// populate order placement results for FinalizeBlock hook
	dextypeswasm.PopulateOrderPlacementResults(contractAddr, dexutils.GetMemState(sdkCtx.Context()).GetAllBlockOrders(sdkCtx, typedContractAddr), cancellations, orderResults)
	dextypeswasm.PopulateOrderExecutionResults(contractAddr, settlements, orderResults)

	return orderResults, settlements, nil
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
