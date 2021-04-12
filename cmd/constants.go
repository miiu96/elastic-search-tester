package cmd

import (
	"github.com/ElrondNetwork/elastic-search-tester/tester"
	"github.com/urfave/cli"
)

var (
	getTxByHash = cli.StringFlag{
		Name:  tester.TxByHash,
		Usage: "The transaction hash",
		Value: "",
	}

	pathToConfig = cli.StringFlag{
		Name:  tester.ConfigFile,
		Usage: "Path to the configuration file",
		Value: "../../config.json",
	}
	verifyBlocks = cli.StringFlag{
		Name:  tester.VerifyAllBlockFromShard,
		Usage: "Will verify all blocks from a given shard",
		Value: "",
	}
)

// GetFlags will return the existing flags
func GetFlags() []cli.Flag {
	return []cli.Flag{
		getTxByHash,
		pathToConfig,
		verifyBlocks,
	}
}
