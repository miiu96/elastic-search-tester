package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/ElrondNetwork/elastic-search-tester/types"
	"github.com/ElrondNetwork/elrond-go/core"
)

const (
	step                     = uint64(1000)
	numberOfBlocksInParallel = 50
)

var (
	errCannotGetBlocks    = errors.New("cannot get blocks")
	errCannotGetMiniblock = errors.New("cannot get miniblocks")
)

type blocksProcessor struct {
	txsProc       TransactionHandler
	elasticClient ElasticHandler
}

func NewBlockProcessor(elasticClient ElasticHandler, txsProc TransactionHandler) *blocksProcessor {
	return &blocksProcessor{
		elasticClient: elasticClient,
		txsProc:       txsProc,
	}
}

func (bp *blocksProcessor) VerifyBlocksMiniblocksAndTransactions(shard string) {
	shardID := uint32(0)
	if shard == "meta" {
		shardID = core.MetachainShardId
	} else {
		s, _ := strconv.ParseUint(shard, 10, 32)
		shardID = uint32(s)
	}

	latestNonce, err := bp.getHighestBlockNonce(shardID)
	log.Warn("highest block nonce", "nonce", latestNonce)
	if err != nil {
		log.Warn("blocksProcessor.getHighestBlockNonce", "error", err)
	}

	for idx := uint64(0); idx < latestNonce; idx += step {
		log.Info("verify blocks...", "from", idx, "to", idx+step)

		encodedQuery := types.Encode(blocksFromTo(shardID, idx, idx+step))

		response, errDo := bp.elasticClient.DoGetRequest(encodedQuery, blocksIndex, "nonce:asc")
		if errDo != nil {
			log.Warn("blocksProcessor.DoGetRequest", "error", err)
			continue
		}

		blocks, errGet := getBlocksFromResponse(response)
		if errGet != nil {
			log.Warn("blocksProcessor.getBlocksFromResponse", "error", err)
			continue
		}

		bp.verifyBlocks(blocks)
	}
}

func (bp *blocksProcessor) verifyBlocks(blocks []*data.Block) {
	wg := &sync.WaitGroup{}
	myChan := make(chan struct{}, numberOfBlocksInParallel)

	for _, block := range blocks {
		wg.Add(1)
		go func(b *data.Block, w *sync.WaitGroup, m chan struct{}) {
			defer func() {
				wg.Done()
				<-m
			}()

			log.Info("block with", "nonce", b.Nonce, "number of miniblocks", len(b.MiniBlocksHashes))

			if len(b.MiniBlocksHashes) == 0 {
				return
			}
			bp.checkMiniblocks(b.MiniBlocksHashes)
			bp.getTransactionsByMBHashes(b.MiniBlocksHashes)

		}(block, wg, myChan)

		myChan <- struct{}{}
	}

	wg.Wait()
}

func (bp *blocksProcessor) getHighestBlockNonce(shardID uint32) (uint64, error) {
	encodedQuery := types.Encode(highestBlock(shardID))

	response, err := bp.elasticClient.DoGetRequest(encodedQuery, blocksIndex, "nonce:desc")
	if err != nil {
		log.Warn("blocksProcessor.DoGetRequest", "error", err)
	}

	blocks, err := getBlocksFromResponse(response)
	if err != nil {
		return 0, err
	}

	return blocks[0].Nonce, nil
}

func (bp *blocksProcessor) getTransactionsByMBHashes(mbsHashes []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(mbsHashes))
	for _, mbHash := range mbsHashes {
		go func(mb string, w *sync.WaitGroup) {
			txs, err := bp.txsProc.GetTransactionsByMBHash(mb)
			if err != nil {
				log.Warn("blocksProcessor.GetTransactionsByMBHash", "error", err)
			}

			bp.txsProc.ParseTransactions(txs)

			w.Done()
		}(mbHash, wg)
	}

	wg.Wait()
}

func (bp *blocksProcessor) checkMiniblocks(mbsHashes []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(mbsHashes))

	for _, mbHash := range mbsHashes {
		go func(mh string, w *sync.WaitGroup) {
			err := bp.checkMiniblock(mh)
			log.LogIfError(err)

			w.Done()
		}(mbHash, wg)
	}

	wg.Wait()
}

func (bp *blocksProcessor) checkMiniblock(mbHash string) error {
	encodedQuery := types.Encode(objByHash(mbHash))
	response, err := bp.elasticClient.DoGetRequest(encodedQuery, miniblocksIndex)
	if err != nil {
		log.Warn("blocksProcessor.checkMiniblock", "error", err)
	}

	hits, ok := response["hits"].(types.ObjectMap)
	if !ok {
		return errCannotGetMiniblock
	}

	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return errCannotGetMiniblock
	}

	sourceElement, ok := source[0].(types.ObjectMap)["_source"]
	if !ok {
		return errCannotGetMiniblock
	}

	sourceBytes, _ := json.Marshal(sourceElement)

	var mb data.Miniblock
	_ = json.Unmarshal(sourceBytes, &mb)

	if mb.ReceiverBlockHash == "" {
		return fmt.Errorf("miniblock hash %s, receiver block hash is empty", mbHash)
	}
	if mb.SenderBlockHash == "" {
		return fmt.Errorf("miniblock hash %s, sender block hash is empty", mbHash)
	}

	return nil
}

func getBlocksFromResponse(response types.ObjectMap) ([]*data.Block, error) {
	hits, ok := response["hits"].(types.ObjectMap)
	if !ok {
		return nil, errCannotGetBlocks
	}

	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return nil, errCannotGetBlocks
	}

	blocks := make([]*data.Block, 0)
	for _, s := range source {
		sourceElement, okS := s.(types.ObjectMap)["_source"]
		if !okS {
			continue
		}

		sourceBytes, _ := json.Marshal(sourceElement)
		var scr data.Block
		_ = json.Unmarshal(sourceBytes, &scr)

		blocks = append(blocks, &scr)
	}

	return blocks, nil
}
