package tests

import (
	"testing"

	"github.com/tolstovrob/graph-go/graph"
)

func TestMakeNode(t *testing.T) {
	node := graph.MakeNode(123)
	if node.Key != 123 {
		t.Errorf("expected key 123, got %d", node.Key)
	}
	if node.Label != "" {
		t.Errorf("expected empty label, got %s", node.Label)
	}

	labeledNode := graph.MakeNode(
		456,
		graph.WithNodeLabel("TestLabel"),
	)
	if labeledNode.Key != 456 {
		t.Errorf("expected key 456, got %d", labeledNode.Key)
	}
	if labeledNode.Label != "TestLabel" {
		t.Errorf("expected label 'TestLabel', got %s", labeledNode.Label)
	}
}

func TestUpdateNode(t *testing.T) {
	node := graph.MakeNode(789)

	node.UpdateNode(graph.WithNodeLabel("UpdatedLabel"))
	if node.Label != "UpdatedLabel" {
		t.Errorf("expected updated label 'UpdatedLabel', got %s", node.Label)
	}

	node.UpdateNode(graph.WithNodeLabel(""))
	if node.Label != "" {
		t.Errorf("expected empty label after update, got %s", node.Label)
	}
}

func TestWithNodeLabel(t *testing.T) {
	node := graph.MakeNode(1)
	opt := graph.WithNodeLabel("LabelForTest")
	opt(node)
	if node.Label != "LabelForTest" {
		t.Errorf("WithNodeLabel option failed, expected 'LabelForTest', got %s", node.Label)
	}
}

func TestMakeEdge(t *testing.T) {
	srcKey := graph.TKey(1)
	dstKey := graph.TKey(2)

	edge := graph.MakeEdge(100, srcKey, dstKey)
	if edge.Key != 100 {
		t.Errorf("expected key 100, got %d", edge.Key)
	}
	if edge.Source != srcKey {
		t.Errorf("expected source %d, got %d", srcKey, edge.Source)
	}
	if edge.Destination != dstKey {
		t.Errorf("expected destination %d, got %d", dstKey, edge.Destination)
	}
	if edge.Weight != 0 {
		t.Errorf("expected default weight 0, got %d", edge.Weight)
	}
	if edge.Label != "" {
		t.Errorf("expected empty label, got %s", edge.Label)
	}

	edgeWithOpts := graph.MakeEdge(
		200,
		srcKey,
		dstKey,
		graph.WithEdgeLabel("Highway"),
		graph.WithEdgeWeight(42),
	)
	if edgeWithOpts.Label != "Highway" {
		t.Errorf("expected label 'Highway', got %s", edgeWithOpts.Label)
	}
	if edgeWithOpts.Weight != 42 {
		t.Errorf("expected weight 42, got %d", edgeWithOpts.Weight)
	}
}

func TestUpdateEdge(t *testing.T) {
	edge := graph.MakeEdge(300, 3, 4)

	edge.UpdateEdge(graph.WithEdgeLabel("Street"), graph.WithEdgeWeight(24))
	if edge.Label != "Street" {
		t.Errorf("expected label 'Street', got %s", edge.Label)
	}
	if edge.Weight != 24 {
		t.Errorf("expected weight 24, got %d", edge.Weight)
	}

	edge.UpdateEdge(graph.WithEdgeLabel(""), graph.WithEdgeWeight(0))
	if edge.Label != "" {
		t.Errorf("expected empty label, got %s", edge.Label)
	}
	if edge.Weight != 0 {
		t.Errorf("expected weight 0, got %d", edge.Weight)
	}
}

func TestWithEdgeLabel(t *testing.T) {
	edge := graph.MakeEdge(1, 1, 2)
	opt := graph.WithEdgeLabel("MainRoad")
	opt(edge)
	if edge.Label != "MainRoad" {
		t.Errorf("WithEdgeLabel option failed, expected 'MainRoad', got %s", edge.Label)
	}
}

func TestWithEdgeWeight(t *testing.T) {
	edge := graph.MakeEdge(1, 1, 2)
	opt := graph.WithEdgeWeight(99)
	opt(edge)
	if edge.Weight != 99 {
		t.Errorf("WithEdgeWeight option failed, expected 99, got %d", edge.Weight)
	}
}

func TestMakeGraphWithOptions(t *testing.T) {
	nodes := []*graph.Node{
		graph.MakeNode(1),
		graph.MakeNode(2, graph.WithNodeLabel("Node2")),
	}

	edges := []*graph.Edge{
		graph.MakeEdge(1, 1, 2),
		graph.MakeEdge(
			2,
			2,
			1,
			graph.WithEdgeWeight(10),
			graph.WithEdgeLabel("Back"),
		),
	}

	options := graph.TOptions{IsMulti: true, IsDirected: true}

	gr := graph.MakeGraph(
		graph.WithGraphNodes(nodes),
		graph.WithGraphEdges(edges),
		graph.WithGraphOptions(options),
	)

	if len(gr.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(gr.Nodes))
	}

	if len(gr.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(gr.Edges))
	}

	if !gr.Options.IsMulti {
		t.Errorf("expected IsMulti true, got false")
	}

	if !gr.Options.IsDirected {
		t.Errorf("expected IsDirected true, got false")
	}
}

func TestUpdateGraphOptions(t *testing.T) {
	gr := graph.MakeGraph()

	gr.UpdateGraph(graph.WithGraphMulti(true))
	if !gr.Options.IsMulti {
		t.Errorf("expected IsMulti true after update, got false")
	}

	gr.UpdateGraph(graph.WithGraphDirected(true))
	if !gr.Options.IsDirected {
		t.Errorf("expected IsDirected true after update, got false")
	}
}

func TestUpdateGraphNodesEdges(t *testing.T) {
	node := graph.MakeNode(5)
	edge := graph.MakeEdge(7, 5, 5)

	gr := graph.MakeGraph()

	gr.UpdateGraph(graph.WithGraphNodes([]*graph.Node{node}))
	if len(gr.Nodes) != 1 || gr.Nodes[0].Key != 5 {
		t.Errorf("expected 1 node with key 5, got %v", gr.Nodes)
	}

	gr.UpdateGraph(graph.WithGraphEdges([]*graph.Edge{edge}))
	if len(gr.Edges) != 1 || gr.Edges[0].Key != 7 {
		t.Errorf("expected 1 edge with key 7, got %v", gr.Edges)
	}
}
