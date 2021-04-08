package tester

import (
	"github.com/ElrondNetwork/elastic-search-tester/client"
	"github.com/ElrondNetwork/elastic-search-tester/config"
	"github.com/ElrondNetwork/elastic-search-tester/process"
	"github.com/urfave/cli"
)

func RunElasticTester(ctx *cli.Context) error {
	pathToConfig := ctx.GlobalString(ConfigFile)
	txHash := ctx.GlobalString(TxByHash)
	verifyBlocksShard := ctx.GlobalString(VerifyAllBlockFromShard)

	cfg, err := config.GetConfig(pathToConfig)
	if err != nil {
		return err
	}

	esClient, err := client.NewElasticClient(cfg)
	if err != nil {
		return err
	}

	txsProc := process.NewTransactionsProc(esClient)
	if txHash != "" {
		txsProc.GetTransactionByHashAndDisplay(txHash)
		return nil
	}
	if verifyBlocksShard != "" {
		blocksProc := process.NewBlockProcessor(esClient, txsProc)
		blocksProc.VerifyBlocksMiniblocksAndTransactions(verifyBlocksShard)
		return nil
	}

	return nil
}
