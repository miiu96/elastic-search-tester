package main

import (
	"os"

	"github.com/ElrondNetwork/elastic-search-tester/cmd"
	"github.com/ElrondNetwork/elastic-search-tester/tester"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("main")

func main() {
	app := cli.NewApp()

	app.Name = "Elastic search tester"
	app.Version = "v1.0.0"
	app.Usage = "This is the entry point for starting elastic search tester"
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}
	app.Flags = cmd.GetFlags()
	app.Action = tester.RunElasticTester

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
