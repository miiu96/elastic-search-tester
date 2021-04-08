package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
)

func TestReflectEtc(t *testing.T) {
	txs := make([]*data.Transaction, 0)
	test(txs)
}

func test(obj interface{}) {
	myType := reflect.ValueOf(obj)

	fmt.Println(reflect.TypeOf(obj).Elem())
	newArray := reflect.SliceOf(reflect.TypeOf(obj).Elem())
	fmt.Println(reflect.TypeOf(newArray))

	fmt.Println(myType)
	fmt.Println(myType.Kind() == reflect.Slice)
}
