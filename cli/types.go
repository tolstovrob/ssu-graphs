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
 * CLI struct represents application state and configuration. It has graph
 * field, which contains info about worked graph. Also it has app fields for
 * configuration of TUI
 */

type CLIService struct {
	app        *tview.Application
	pages      *tview.Pages
	statusView *tview.TextView
	graph      *graph.Graph
}

/*
 * Status type which essential for status bar
 */

type Status int

const (
	Ok Status = iota
	Error
)

/*
 * This is not neccesarily a type, but since it reffered to status, I decided
 * to put it here in types.go
 */

var statusColor = map[Status]string{
	Ok:    "white",
	Error: "red",
}
