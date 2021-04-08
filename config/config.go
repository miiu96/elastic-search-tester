package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("main")

type Config struct {
	ElasticURL string `json:"elastic-url"`
	User       string `json:"user"`
	Password   string `json:"password"`
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("readFile: open file %s error %s", fileName, err)
	}
	defer func() {
		errClose := file.Close()
		log.LogIfError(errClose)
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("readFile: ioutil.ReadAll %s error %s", fileName, err)
	}
	return data, nil
}

func GetConfig(path string) (*Config, error) {
	fileContent, err := readFile(path)
	if err != nil {
		return nil, err
	}
	servConfig := &Config{}
	err = json.Unmarshal(fileContent, servConfig)
	if err != nil {
		return nil, err
	}

	return servConfig, nil
}
