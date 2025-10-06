package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tolstovrob/graph-go/graph"
)

func RunCLI() {
	reader := bufio.NewReader(os.Stdin)
	gr := graph.MakeGraph()

	fmt.Println("Graph CLI interface. Type 'help' for commands.")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Split(line, " ")
		switch args[0] {
		case "help":
			printHelp()
		case "exit":
			return
		case "addnode":
			addNode(gr, args[1:])
		case "removenode":
			removeNode(gr, args[1:])
		case "addedge":
			addEdge(gr, args[1:])
		case "removeedge":
			removeEdge(gr, args[1:])
		case "listnodes":
			listNodes(gr)
		case "listedges":
			listEdges(gr)
		case "listadj":
			listAdjacency(gr)
		case "save":
			saveGraph(gr, args[1:])
		case "load":
			loadGraph(gr, args[1:])
		default:
			fmt.Println("Unknown command. Type 'help' for list of commands.")
		}
	}
}

func printHelp() {
	fmt.Println(`Commands:
  help - show this help message
  exit - exit the CLI
  addnode <key> [label] - add a node with given key and optional label
  removenode <key> - remove node by key
  addedge <key> <source> <destination> [weight] [label] - add edge, weight and label optional
  removeedge <key> - remove edge by key
  listnodes - list all nodes
  listedges - list all edges
  listadj - display adjacency map
  save <filename> - save graph to file
  load <filename> - load graph from file
`)
}

func addNode(gr *graph.Graph, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: addnode <key> [label]")
		return
	}
	key, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid key:", args[0])
		return
	}
	var label string
	if len(args) > 1 {
		label = strings.Join(args[1:], " ")
	}
	node := graph.MakeNode(graph.TKey(key), graph.WithNodeLabel(label))
	err = gr.AddNode(node)
	if err != nil {
		fmt.Println("Error adding node:", err)
		return
	}
	fmt.Println("Node added.")
}

func removeNode(gr *graph.Graph, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: removenode <key>")
		return
	}
	key, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid key:", args[0])
		return
	}
	err = gr.RemoveNodeByKey(graph.TKey(key))
	if err != nil {
		fmt.Println("Error removing node:", err)
		return
	}
	fmt.Println("Node removed.")
}

func addEdge(gr *graph.Graph, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: addedge <key> <source> <destination> [weight] [label]")
		return
	}
	key, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid key:", args[0])
		return
	}
	src, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		fmt.Println("Invalid source:", args[1])
		return
	}
	dst, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		fmt.Println("Invalid destination:", args[2])
		return
	}

	var weight graph.TWeight = 0
	var label string
	if len(args) > 3 {
		w, err := strconv.ParseUint(args[3], 10, 64)
		if err == nil {
			weight = graph.TWeight(w)
			if len(args) > 4 {
				label = strings.Join(args[4:], " ")
			}
		} else {
			label = strings.Join(args[3:], " ")
		}
	}

	edge := graph.MakeEdge(graph.TKey(key), graph.TKey(src), graph.TKey(dst), graph.WithEdgeWeight(weight), graph.WithEdgeLabel(label))
	err = gr.AddEdge(edge)
	if err != nil {
		fmt.Println("Error adding edge:", err)
		return
	}
	fmt.Println("Edge added.")
}

func removeEdge(gr *graph.Graph, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: removeedge <key>")
		return
	}
	key, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid key:", args[0])
		return
	}
	err = gr.RemoveEdgeByKey(graph.TKey(key))
	if err != nil {
		fmt.Println("Error removing edge:", err)
		return
	}
	fmt.Println("Edge removed.")
}

func listNodes(gr *graph.Graph) {
	if len(gr.Nodes) == 0 {
		fmt.Println("No nodes in graph.")
		return
	}
	fmt.Println("Nodes:")
	for k, node := range gr.Nodes {
		fmt.Printf("Key: %d, Label: %s\n", k, node.Label)
	}
}

func listEdges(gr *graph.Graph) {
	if len(gr.Edges) == 0 {
		fmt.Println("No edges in graph.")
		return
	}
	fmt.Println("Edges:")
	for k, edge := range gr.Edges {
		fmt.Printf("Key: %d, From: %d, To: %d, Weight: %d, Label: %s\n", k, edge.Source, edge.Destination, edge.Weight, edge.Label)
	}
}

func listAdjacency(gr *graph.Graph) {
	if len(gr.AdjacencyMap) == 0 {
		fmt.Println("Adjacency map is empty.")
		return
	}
	fmt.Println("Adjacency:")
	for key, neighbors := range gr.AdjacencyMap {
		fmt.Printf("%d: %v\n", key, neighbors)
	}
}

func saveGraph(gr *graph.Graph, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: save <filename>")
		return
	}
	filename := args[0]
	jsonStr, err := gr.ToJSON()
	if err != nil {
		fmt.Println("Error serializing graph:", err)
		return
	}
	err = os.WriteFile(filename, []byte(jsonStr), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	fmt.Println("Graph saved to", filename)
}

func loadGraph(gr *graph.Graph, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: load <filename>")
		return
	}
	filename := args[0]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	err = gr.FromJSON(string(data))
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}
	fmt.Println("Graph loaded from", filename)
}
