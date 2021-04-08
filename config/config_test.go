package config

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Parallel()

	cfg, err := GetConfig("../config.json")

	fmt.Println(err)
	fmt.Println(cfg)
}
