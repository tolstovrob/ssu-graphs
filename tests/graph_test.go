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

	labeledNode := graph.MakeNode(456, graph.WithNodeLabel("TestLabel"))
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
