package process

import (
	"bytes"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/ElrondNetwork/elastic-search-tester/types"
)

type ElasticHandler interface {
	DoGetRequest(query *bytes.Buffer, index string, sort ...string) (types.ObjectMap, error)
}

type TransactionHandler interface {
	GetTransactionsByMBHash(mbHash string) ([]*data.Transaction, error)
	ParseTransactions(txs []*data.Transaction)
}
