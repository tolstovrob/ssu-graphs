/*
 * This a CLI service for my graph implementation. It is build with tview and
 * represents TUI CLI.
 *
 * Author: github.com/tolstovrob
 */

package cli

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (cli *CLIService) setupUI() {
	cli.pages = tview.NewPages()

	// create and push main page as 'main'. It will be referenced as main
	mainMenu := cli.createMainMenu()
	cli.pages.AddPage("main", mainMenu, true, true)

	// create status line
	cli.statusView = tview.
		NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	cli.statusView.
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetBorder(true)

	// arrange status line and main page
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cli.pages, 0, 1, true).
		AddItem(cli.statusView, 3, 0, false)

	cli.app.SetRoot(flex, true)
	cli.updateStatus("Ready. Use arrows or tab to navigate, q to exit", Default)
}

func (cli *CLIService) createMainMenu() tview.Primitive {
	list := tview.NewList().
		AddItem("Node Operations", "Add, remove, modify nodes", '1', cli.showNodeOperations).
		AddItem("Edge Operations", "Add, remove, modify edges", '2', cli.showEdgeOperations).
		AddItem("Graph Options", "Configure graph properties", '3', cli.showGraphOptions).
		AddItem("View Graph Info", "Display graph information", '4', cli.showGraphInfo).
		AddItem("JSON Operations", "Save/Load graph from JSON", '5', cli.showJSONOperations).
		AddItem("Algorithms", "Tasks from my SSU course", '6', cli.showAlgorithmsMenu).
		AddItem("Quit", "Exit application", 'q', func() {
			cli.app.Stop()
		})

	list.SetBorder(true).SetTitle(" Graph CLI - Main Menu ")
	return list
}

func (cli *CLIService) updateStatus(message string, status Status) {
	cli.statusView.SetText(fmt.Sprintf("[%s]%s", statusColor[status], message))
}

func (cli *CLIService) showScrollableModal(title, content, returnPage string) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		}).
		SetText(content)

	textView.SetBorder(true).SetTitle(fmt.Sprintf(" %s - Use keys to scroll, Q to go back ", title))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			cli.pages.SwitchToPage(returnPage)
			return nil
		}
		return event
	})

	cli.pages.AddAndSwitchToPage(strings.ToLower(strings.ReplaceAll(title, " ", "_"))+"_view", flex, true)
}
