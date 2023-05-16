package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/app/params"
	"github.com/spf13/cobra"
	leveldbutils "github.com/syndtr/goleveldb/leveldb/util"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tm-db"
)

func CompactCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compact",
		Short: "Compact the application DB fully (only if it is a levelDB)",
		Long:  `Compact the application DB fully (only if it is a levelDB)`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			params.SetTendermintConfigs(config)
			config.SetRoot(clientCtx.HomeDir)
			rootDir := config.RootDir
			dataDir := filepath.Join(rootDir, "data")
			db, err := sdk.NewLevelDB("application", dataDir)
			if err != nil {
				return err
			}
			if goleveldb, ok := db.(*dbm.GoLevelDB); ok {
				start := time.Now()
				if err := goleveldb.DB().CompactRange(leveldbutils.Range{Start: nil, Limit: nil}); err != nil {
					return err
				}
				fmt.Printf("compaction took %f seconds\n", time.Since(start).Seconds())

				keyCnt, keySize, valSize := 0, 0, 0
				iter, err := goleveldb.Iterator(nil, nil)
				if err != nil {
					return err
				}
				for ; iter.Valid(); iter.Next() {
					keyCnt++
					if keyCnt%10000 == 0 {
						fmt.Printf("Iterated through %d KV pairs\n", keyCnt)
					}
					keySize += len(iter.Key())
					valSize += len(iter.Value())
				}
				fmt.Printf("KV count: %d, total key size: %d, total value size: %d\n", keyCnt, keySize, valSize)
			} else {
				return errors.New("cannot compact non-levelDB")
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")

	return cmd
}
