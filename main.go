package main

import "github.com/tolstovrob/graph-go/cli"

func main() {
	cliService := cli.NewCLIService()
	if err := cliService.Run(); err != nil {
		panic(err)
	}
}
