package process

import (
	"testing"

	"github.com/ElrondNetwork/elastic-search-tester/client"
	"github.com/ElrondNetwork/elastic-search-tester/config"
)

func TestGetTransactionByHash(t *testing.T) {
	elasticC, _ := client.NewElasticClient(&config.Config{
		ElasticURL: "https://search-testing-mihai-cupm4ru4fsbpsgkikuqx6oexie.eu-central-1.es.amazonaws.com",
	})

	txsProc := NewTransactionsProc(elasticC)

	txsProc.GetTransactionByHashAndDisplay("a1f999fd00eeb47cae3448b40030ad75038c5eb25fc3f4d3f74b879e6426bc2a")
}
