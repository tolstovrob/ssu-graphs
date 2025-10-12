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

func (cli *CLIService) showEdgeOperations() {
	modal := tview.NewModal().
		SetText("Edge Operations").
		AddButtons([]string{"Add Edge", "Remove Edge", "Modify Edge", "List Edges", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Add Edge":
				cli.showAddEdgeForm()
			case "Remove Edge":
				cli.showRemoveEdgeForm()
			case "Modify Edge":
				cli.showModifyEdgeForm()
			case "List Edges":
				cli.showEdgesList()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("edge_operations", modal, true)
}
func (cli *CLIService) showAddEdgeForm() {
	form := tview.NewForm()
	var edgeKey, srcKey, dstKey, weightStr, label string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		edgeKey = text
	})
	form.AddInputField("Source Node Key", "", 10, nil, func(text string) {
		srcKey = text
	})
	form.AddInputField("Destination Node Key", "", 10, nil, func(text string) {
		dstKey = text
	})
	form.AddInputField("Weight", "", 10, nil, func(text string) {
		weightStr = text
	})
	form.AddInputField("Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Add", func() {
		key, err := strconv.ParseUint(edgeKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid edge key format", Error)
			return
		}

		src, err := strconv.ParseUint(srcKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid source key format", Error)
			return
		}

		dst, err := strconv.ParseUint(dstKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid destination key format", Error)
			return
		}

		weight, err := strconv.ParseUint(weightStr, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid weight format", Error)
			return
		}

		edge := graph.MakeEdge(graph.TKey(key), graph.TKey(src), graph.TKey(dst))
		if weight > 0 {
			edge.UpdateEdge(graph.WithEdgeWeight(graph.TWeight(weight)))
		}
		if label != "" {
			edge.UpdateEdge(graph.WithEdgeLabel(label))
		}

		if err := cli.graph.AddEdge(edge); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
		} else {
			cli.updateStatus(fmt.Sprintf("Edge %d added successfully", key), Success)
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Add Edge ")
	cli.pages.AddAndSwitchToPage("add_edge", form, true)
}

func (cli *CLIService) showRemoveEdgeForm() {
	form := tview.NewForm()
	var key string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddButton("Remove", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if err := cli.graph.RemoveEdgeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
		} else {
			cli.updateStatus(fmt.Sprintf("Edge %d removed successfully", keyVal), Success)
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Remove Edge ")
	cli.pages.AddAndSwitchToPage("remove_edge", form, true)
}

func (cli *CLIService) showModifyEdgeForm() {
	form := tview.NewForm()
	var key, weightStr, label string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("New Weight", "", 10, nil, func(text string) {
		weightStr = text
	})
	form.AddInputField("New Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Modify", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		edge, err := cli.graph.GetEdgeByKey(graph.TKey(keyVal))
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
			return
		}

		if weightStr != "" {
			weight, err := strconv.ParseUint(weightStr, 10, 64)
			if err != nil {
				cli.updateStatus("Error: Invalid weight format", Error)
				return
			}
			edge.UpdateEdge(graph.WithEdgeWeight(graph.TWeight(weight)))
		}

		if label != "" {
			edge.UpdateEdge(graph.WithEdgeLabel(label))
		}

		cli.updateStatus(fmt.Sprintf("Edge %d modified successfully", keyVal), Success)
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Modify Edge ")
	cli.pages.AddAndSwitchToPage("modify_edge", form, true)
}

func (cli *CLIService) showEdgesList() {
	edgesInfo := "Edges:\n\n"
	for key, edge := range cli.graph.Edges {
		edgesInfo += fmt.Sprintf("Key: %d, Source: %d -> Destination: %d, Weight: %d, Label: %s\n",
			key, edge.Source, edge.Destination, edge.Weight, edge.Label)
	}

	cli.showScrollableModal("Edges List", edgesInfo, "edge_operations")
}
