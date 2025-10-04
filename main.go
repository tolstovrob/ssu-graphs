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

	labeledNode := graph.MakeNode(1, graph.WithNodeLabel("Aboba"))
	fmt.Printf("Created node with key %v and label %v\n", labeledNode.Key, labeledNode.Label)
}
