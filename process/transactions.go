package process

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/ElrondNetwork/elastic-search-tester/types"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var (
	log                     = logger.GetOrCreate("process")
	errCannotGetTransaction = errors.New("cannot get transaction")
	errCannotGetSCRS        = errors.New("cannot get smart contracts results")
	errCannotGetReceipt     = errors.New("cannot get receipts")
)

type transactionsProc struct {
	elasticHandler ElasticHandler
}

func NewTransactionsProc(elasticHandler ElasticHandler) *transactionsProc {
	return &transactionsProc{
		elasticHandler: elasticHandler,
	}
}

func (tp *transactionsProc) GetTransactionByHashAndDisplay(txHash string) {
	encodedQuery := types.Encode(objByHash(txHash))

	response, err := tp.elasticHandler.DoGetRequest(encodedQuery, transactionsIndex)
	if err != nil {
		log.Warn("transactionsProc.GetTransactionByHashAndDisplay", "error", err)
		return
	}

	tx, err := getTxFromResponse(response)
	if err != nil {
		log.Warn("transactionsProc.GetTransactionByHashAndDisplay", "error", err)
		return
	}

	fmt.Println("TRANSACTION")
	prettyPrint(tx)

	tp.GetTxResults(txHash, tx)
}

func prettyPrint(obj interface{}) {
	objBytes, _ := json.Marshal(obj)
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, objBytes, "", "\t")

	fmt.Println(prettyJSON.String())
}

func (tp *transactionsProc) GetTxResults(txHash string, tx *data.Transaction) {
	if tx.HasSCR {
		scrs, err := tp.GetSCR(txHash)
		if err != nil {
			log.Warn("transactionsProc.getAndPrintTxResults", "error", err)
		}

		fmt.Println("SMART CONTRACT RESULTS")
		for _, scr := range scrs {
			prettyPrint(scr)
		}
	}
	if tx.Status == "invalid" {
		rec, err := tp.GetReceipt(txHash)
		if err != nil {
			log.Warn("transactionsProc.getAndPrintTxResults", "error", err)
		}

		fmt.Println("RECEIPT")
		prettyPrint(rec)
		return
	}
}

func (tp *transactionsProc) GetSCR(txHash string) ([]*data.ScResult, error) {
	encodedQuery := types.Encode(scrsByTxHash(txHash))
	response, err := tp.elasticHandler.DoGetRequest(encodedQuery, scrsIndex)
	if err != nil {
		return nil, err
	}

	hits, ok := response["hits"].(types.ObjectMap)
	if !ok {
		return nil, errCannotGetSCRS
	}

	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return nil, errCannotGetSCRS
	}

	scrs := make([]*data.ScResult, 0)
	for _, s := range source {
		sourceElement, okS := s.(types.ObjectMap)["_source"]
		if !okS {
			continue
		}

		sourceBytes, _ := json.Marshal(sourceElement)
		var scr data.ScResult
		_ = json.Unmarshal(sourceBytes, &scr)

		scrs = append(scrs, &scr)
	}

	return scrs, nil
}

func (tp *transactionsProc) GetReceipt(txHash string) (*data.Receipt, error) {
	encodedQuery := types.Encode(receiptByTxHash(txHash))
	response, err := tp.elasticHandler.DoGetRequest(encodedQuery, receiptsIndex)
	if err != nil {
		return nil, err
	}

	hits, ok := response["hits"].(types.ObjectMap)
	if !ok {
		return nil, errCannotGetReceipt
	}

	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return nil, errCannotGetReceipt
	}

	sourceElement, ok := source[0].(types.ObjectMap)["_source"]
	if !ok {
		return nil, errCannotGetReceipt
	}

	sourceBytes, _ := json.Marshal(sourceElement)

	var rec data.Receipt
	_ = json.Unmarshal(sourceBytes, &rec)

	return &rec, nil
}

func getTxFromResponse(res types.ObjectMap) (*data.Transaction, error) {
	hits, ok := res["hits"].(types.ObjectMap)
	if !ok {
		return nil, errCannotGetTransaction
	}

	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return nil, errCannotGetTransaction
	}

	sourceElement, ok := source[0].(types.ObjectMap)["_source"]
	if !ok {
		return nil, errCannotGetTransaction
	}

	sourceBytes, _ := json.Marshal(sourceElement)

	var tx data.Transaction
	_ = json.Unmarshal(sourceBytes, &tx)

	return &tx, nil
}

func (tp *transactionsProc) GetTransactionsByMBHash(mbHash string) ([]*data.Transaction, error) {
	encodedQuery := types.Encode(txsByMBHash(mbHash))

	response, err := tp.elasticHandler.DoGetRequest(encodedQuery, transactionsIndex)
	if err != nil {
		return nil, err
	}

	hits, ok := response["hits"].(types.ObjectMap)
	if !ok {
		return nil, errCannotGetTransaction
	}

	txs := make([]*data.Transaction, 0)
	source, ok := hits["hits"].([]interface{})
	if !ok || len(source) < 1 {
		return txs, nil
	}

	for _, s := range source {
		ss := s.(types.ObjectMap)

		sourceElement, okS := ss["_source"]
		if !okS {
			continue
		}

		sourceBytes, _ := json.Marshal(sourceElement)
		var tx data.Transaction
		_ = json.Unmarshal(sourceBytes, &tx)

		tx.Hash = fmt.Sprintf("%v", ss["_id"])

		txs = append(txs, &tx)
	}

	return txs, nil
}

func (tp *transactionsProc) ParseTransactions(txs []*data.Transaction) {
	wg := &sync.WaitGroup{}
	myChan := make(chan struct{}, 150)

	log.Info("number of transactions", "num", len(txs))
	for _, tx := range txs {
		wg.Add(1)
		go func(t *data.Transaction, w *sync.WaitGroup, m chan struct{}) {
			defer func() {
				w.Done()
				<-m
			}()

			if t.HasSCR == true {
				scrs, err := tp.GetSCR(t.Hash)
				if err != nil {
					log.Warn("cannot get smart contract results for tx", "hash", t.Hash, "error", err)
					return
				}
				if len(scrs) == 0 {
					log.Warn("transactions should have SCRS", "hash", t.Hash, "error", err)
					return
				}

				return
			}

			if t.Status == "invalid" {
				_, err := tp.GetReceipt(t.Hash)
				if err != nil {
					log.Warn("cannot get receipt for tx", "hash", t.Hash, "error", err)
					return
				}

				return
			}

			if t.Status == "pending" {
				log.Warn("pending transaction", "hash", t.Hash, "timestamp", time.Unix(int64(t.Timestamp), 0).Format(time.RFC822Z))
			}

		}(tx, wg, myChan)

		myChan <- struct{}{}
	}

	wg.Wait()
}
