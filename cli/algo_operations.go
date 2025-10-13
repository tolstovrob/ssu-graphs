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
	"github.com/tolstovrob/graph-go/algo"
	"github.com/tolstovrob/graph-go/graph"
)

func (cli *CLIService) showAlgorithmsMenu() {
	list := tview.NewList().
		AddItem("In-Degree less than", "Find nodes with in-degree less than target", '1', cli.showInDegreeLessThanForm).
		AddItem("In-nodes in directed", "Find nodes, that are in-nodes for target in directed graph", '2', cli.showIncomingNeighborsForm).
		AddItem("Back to Main Menu", "Return to main menu", 'q', func() {
			cli.pages.SwitchToPage("main")
		})

	list.SetBorder(true).SetTitle(" Graph Algorithms ")
	cli.pages.AddAndSwitchToPage("algorithms_menu", list, true)
}

func (cli *CLIService) showInDegreeLessThanForm() {
	form := tview.NewForm()
	var targetKey string

	form.AddInputField("Target Node Key", "", 10, nil, func(text string) {
		targetKey = text
	})
	form.AddButton("Run Algorithm", func() {
		keyVal, err := strconv.ParseUint(targetKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if _, err := cli.graph.GetNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: Node %d does not exist", keyVal), Error)
			return
		}

		result := algo.InDegreeLessThan(cli.graph, graph.TKey(keyVal))

		var resultText string
		if len(result) == 0 {
			resultText = fmt.Sprintf("No nodes found with in-degree less than target node %d", keyVal)
		} else {
			resultText = fmt.Sprintf("Nodes with in-degree less than target node %d:\n\n", keyVal)
			for i, nodeKey := range result {
				node, _ := cli.graph.GetNodeByKey(nodeKey)
				if node != nil && node.Label != "" {
					resultText += fmt.Sprintf("%d. Node %d (Label: %s)\n", i+1, nodeKey, node.Label)
				} else {
					resultText += fmt.Sprintf("%d. Node %d\n", i+1, nodeKey)
				}
			}
			resultText += fmt.Sprintf("\nTotal: %d nodes", len(result))
		}

		cli.showScrollableModal("Algorithm Result", resultText, "algorithms_menu")
		cli.updateStatus("Algorithm completed successfully", Success)
	})

	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("algorithms_menu")
	})

	form.SetBorder(true).SetTitle(" In-Degree Less Than Algorithm ")
	cli.pages.AddAndSwitchToPage("in_degree_algorithm", form, true)
}

func (cli *CLIService) showIncomingNeighborsForm() {
	form := tview.NewForm()
	var targetKey string

	form.AddInputField("Target Vertex Key", "", 10, nil, func(text string) {
		targetKey = text
	})
	form.AddButton("Find Neighbors", func() {
		keyVal, err := strconv.ParseUint(targetKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format", Error)
			return
		}

		if _, err := cli.graph.GetNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: Node %d does not exist", keyVal), Error)
			return
		}

		result, err := algo.InNodesInDirected(cli.graph, graph.TKey(keyVal))

		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err), Error)
			return
		}

		var resultText string

		resultText = fmt.Sprintf("Incoming neighbors for vertex %d (directed graph):\n\n", keyVal)

		if len(result) == 0 {
			resultText += "No incoming neighbors found"
		} else {
			for i, neighborKey := range result {
				node, _ := cli.graph.GetNodeByKey(neighborKey)
				if node != nil && node.Label != "" {
					resultText += fmt.Sprintf("%d. Node %d (Label: %s)\n", i+1, neighborKey, node.Label)
				} else {
					resultText += fmt.Sprintf("%d. Node %d\n", i+1, neighborKey)
				}
			}
			resultText += fmt.Sprintf("\nTotal: %d incoming neighbors", len(result))
		}

		cli.showScrollableModal("Incoming Neighbors", resultText, "algorithms_menu")
		cli.updateStatus("Incoming neighbors found successfully", Success)
	})

	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("algorithms_menu")
	})

	form.SetBorder(true).SetTitle(" Find Incoming Neighbors ")
	cli.pages.AddAndSwitchToPage("incoming_neighbors", form, true)
}
