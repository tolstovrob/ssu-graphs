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
