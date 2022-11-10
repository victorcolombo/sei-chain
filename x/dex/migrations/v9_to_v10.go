package migrations

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/x/dex/keeper"
	"github.com/sei-protocol/sei-chain/x/dex/types"
	"github.com/tendermint/tendermint/libs/log"
)

func V9ToV10(ctx sdk.Context, dexkeeper keeper.Keeper, logger log.Logger) error {
	allContractInfo := dexkeeper.GetAllContractInfo(ctx)
	for _, contractInfo := range allContractInfo {

		store := prefix.NewStore(
			ctx.KVStore(dexkeeper.StoreKey),
			types.MatchResultPrefix(contractInfo.ContractAddr),
		)
		prevHeight := 9736346
		// Get latest match result
		key := make([]byte, 8)
		//binary.BigEndian.PutUint64(key, uint64(1))
		//logger.Error("Match result key exists for 1 %t", store.Has(key))
		//binary.BigEndian.PutUint64(key, uint64(9736346))
		//logger.Error("Match result key exists for 9736346 %t", store.Has(key))

		binary.BigEndian.PutUint64(key, uint64(prevHeight))
		if !store.Has(key) {
			panic(fmt.Sprintf("Match result key not found for height %d", prevHeight))
		}
		bz := store.Get(key)
		result := types.MatchResult{}
		if err := result.Unmarshal(bz); err != nil {
			panic(err)
		}
		dexkeeper.SetMatchResult(ctx, contractInfo.ContractAddr, &result)

		// Now, remove all older ones
		for i := 0; i <= prevHeight; i++ {
			key := make([]byte, 8)
			binary.BigEndian.PutUint64(key, uint64(i))
			store.Delete(key)
		}
	}
	return nil
}
