package main

import (
	"fmt"

	"github.com/tolstovrob/graph-go/graph"
)

func main() {
	fmt.Println("Yet another graph implementation")

	node := graph.MakeNode(1)
	fmt.Printf("Created node with key %v and label %v\n", node.Key, node.Label)

	node.UpdateNode(graph.WithNodeLabel("Fimoz"))
	fmt.Printf("Updated label for node with key %v: %v\n", node.Key, node.Label)

	labeledNode := graph.MakeNode(2, graph.WithNodeLabel("Aboba"))
	fmt.Printf("Created node with key %v and label %v\n", labeledNode.Key, labeledNode.Label)

	edge := graph.MakeEdge(1, node.Key, labeledNode.Key)
	fmt.Printf("Connected %v with %v. Connection key: %v\n", node.Key, labeledNode.Key, edge.Key)
	fmt.Printf("Label edge %v with %v. Weight it as %v\n", edge.Key, edge.Label, edge.Weight)

	edge.UpdateEdge(graph.WithEdgeLabel("Zalupa"), graph.WithEdgeWeight(5))
	fmt.Printf("Rename edge %v with %v. Weight it as %v\n", edge.Key, edge.Label, edge.Weight)

	gr := graph.Graph{}
	dgr := graph.MakeGraph(graph.WithGraphDirected(true))
	fmt.Printf("First graph is directed: %v, and second is directed: %v\n", gr.Options.IsDirected, dgr.Options.IsDirected)
}
