/*
 * This a CLI service for my graph implementation. It is build with tview and
 * represents TUI CLI.
 *
 * Author: github.com/tolstovrob
 */

package cli

import (
	"github.com/rivo/tview"
	"github.com/tolstovrob/graph-go/graph"
)

/*
 * A constructor for CLI
 */

func NewCLIService() *CLIService {
	cli := &CLIService{
		app:   tview.NewApplication(),
		graph: graph.MakeGraph(),
	}

	cli.setupUI()
	return cli
}

/*
 * CLI runner. Just use:
 *
 * cli := NewCLIService()
 * if err := cliService.Run(); err != nil {
 *   // handle your error
 * }
 */

func (cli *CLIService) Run() error {
	return cli.app.Run()
}
