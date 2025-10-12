/*
 * This a CLI service for my graph implementation. It is build with tview and
 * represents TUI CLI.
 *
 * Author: github.com/tolstovrob
 */

package cli

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/tolstovrob/graph-go/graph"
)

func (cli *CLIService) showNodeOperations() {
	modal := tview.NewModal().
		SetText("Node Operations").
		AddButtons([]string{"Add Node", "Remove Node", "Modify Node", "List Nodes", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Add Node":
				cli.showAddNodeForm()
			case "Remove Node":
				cli.showRemoveNodeForm()
			case "Modify Node":
				cli.showModifyNodeForm()
			case "List Nodes":
				cli.showNodesList()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("node_operations", modal, true)
}
func (cli *CLIService) showAddNodeForm() {
	form := tview.NewForm()
	var key, label string

	form.AddInputField("Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Add", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		node := graph.MakeNode(graph.TKey(keyVal))
		if label != "" {
			node.UpdateNode(graph.WithNodeLabel(label))
		}

		if err := cli.graph.AddNode(node); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
		} else {
			cli.updateStatus(fmt.Sprintf("Node %d added successfully", keyVal), Error)
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Add Node ")
	cli.pages.AddAndSwitchToPage("add_node", form, true)
}

func (cli *CLIService) showRemoveNodeForm() {
	form := tview.NewForm()
	var key string

	form.AddInputField("Node Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddButton("Remove", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if err := cli.graph.RemoveNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
		} else {
			cli.updateStatus(fmt.Sprintf("Node %d removed successfully", keyVal), Error)
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Remove Node ")
	cli.pages.AddAndSwitchToPage("remove_node", form, true)
}

func (cli *CLIService) showModifyNodeForm() {
	form := tview.NewForm()
	var key, newLabel string

	form.AddInputField("Node Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("New Label", "", 20, nil, func(text string) {
		newLabel = text
	})
	form.AddButton("Modify", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		node, err := cli.graph.GetNodeByKey(graph.TKey(keyVal))
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
			return
		}

		node.UpdateNode(graph.WithNodeLabel(newLabel))
		cli.updateStatus(fmt.Sprintf("Node %d modified successfully", keyVal), Error)
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Modify Node ")
	cli.pages.AddAndSwitchToPage("modify_node", form, true)
}

func (cli *CLIService) showNodesList() {
	nodesInfo := "Nodes:\n\n"
	for key, node := range cli.graph.Nodes {
		nodesInfo += fmt.Sprintf("Key: %d, Label: %s\n", key, node.Label)
	}

	cli.showScrollableModal("Nodes List", nodesInfo, "node_operations")
}
