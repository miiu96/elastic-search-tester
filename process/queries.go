package process

import "github.com/ElrondNetwork/elastic-search-tester/types"

func objByHash(hash string) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"match": types.ObjectMap{
				"_id": hash,
			},
		},
	}
}

func scrsByTxHash(hash string) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"match": types.ObjectMap{
				"originalTxHash": hash,
			},
		},
	}
}

func receiptByTxHash(hash string) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"match": types.ObjectMap{
				"txHash": hash,
			},
		},
	}
}

func blocksFromTo(shardID uint32, from, to uint64) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"bool": types.ObjectMap{
				"must": []interface{}{
					types.ObjectMap{
						"match": types.ObjectMap{
							"shardId": shardID,
						},
					},
					types.ObjectMap{
						"range": types.ObjectMap{
							"nonce": types.ObjectMap{
								"gte": from,
								"lte": to,
							},
						},
					},
				},
			},
		},
		"size": 1000,
	}
}

func highestBlock(shardID uint32) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"bool": types.ObjectMap{
				"must": []interface{}{
					types.ObjectMap{
						"match": types.ObjectMap{
							"shardId": shardID,
						},
					},
				},
			},
		},
		"size": 1,
	}
}

func txsByMBHash(mbHash string) types.ObjectMap {
	return types.ObjectMap{
		"query": types.ObjectMap{
			"match": types.ObjectMap{
				"miniBlockHash": mbHash,
			},
		},
		"size": 9999,
	}
}
