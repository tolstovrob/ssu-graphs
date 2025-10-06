package graph_test

import (
	"testing"

	"github.com/tolstovrob/graph-go/graph"
)

func TestMakeNode(t *testing.T) {
	node := graph.MakeNode(1)
	if node.Key != 1 {
		t.Errorf("Expected node key to be 1, got %d", node.Key)
	}
}

func TestMakeNodeWithLabel(t *testing.T) {
	label := "Aboba"
	node := graph.MakeNode(2, graph.WithNodeLabel(label))
	if node.Label != label {
		t.Errorf("Expected label %s, got %s", label, node.Label)
	}
}

func TestUpdateNode(t *testing.T) {
	node := graph.MakeNode(3)
	newLabel := "New Label"
	node.UpdateNode(graph.WithNodeLabel(newLabel))
	if node.Label != newLabel {
		t.Errorf("Expected label to be updated to %s, got %s", newLabel, node.Label)
	}
}

func TestMakeEdge(t *testing.T) {
	src := graph.MakeNode(1)
	dst := graph.MakeNode(2)
	edge := graph.MakeEdge(1, src.Key, dst.Key)
	if edge.Key != 1 || edge.Source != src.Key || edge.Destination != dst.Key {
		t.Errorf("Edge fields do not match expected values")
	}
}

func TestMakeEdgeWithOptions(t *testing.T) {
	src := graph.MakeNode(1)
	dst := graph.MakeNode(2)
	weight := graph.TWeight(69)
	label := "Path"
	edge := graph.MakeEdge(2, src.Key, dst.Key, graph.WithEdgeWeight(weight), graph.WithEdgeLabel(label))
	if edge.Weight != weight {
		t.Errorf("Expected weight %d, got %d", weight, edge.Weight)
	}
	if edge.Label != label {
		t.Errorf("Expected label %s, got %s", label, edge.Label)
	}
}

func TestUpdateEdge(t *testing.T) {
	edge := graph.MakeEdge(1, 1, 2)
	newWeight := graph.TWeight(100)
	edge.UpdateEdge(graph.WithEdgeWeight(newWeight))
	if edge.Weight != newWeight {
		t.Errorf("Expected weight to be updated to %d, got %d", newWeight, edge.Weight)
	}
}

func TestAddGetRemoveNode(t *testing.T) {
	gr := graph.MakeGraph()
	node := graph.MakeNode(1)
	if err := gr.AddNode(node); err != nil {
		t.Errorf("Failed to add node: %v", err)
	}
	n, err := gr.GetNodeByKey(node.Key)
	if err != nil || n == nil {
		t.Errorf("Failed to get node: %v", err)
	}
	if err := gr.RemoveNodeByKey(node.Key); err != nil {
		t.Errorf("Failed to remove node: %v", err)
	}
	if _, err := gr.GetNodeByKey(node.Key); err == nil {
		t.Errorf("Expected error getting removed node")
	}
}

func TestAddGetRemoveEdge(t *testing.T) {
	gr := graph.MakeGraph()
	src := graph.MakeNode(1)
	dst := graph.MakeNode(2)
	gr.AddNode(src)
	gr.AddNode(dst)
	edge := graph.MakeEdge(1, src.Key, dst.Key)
	if err := gr.AddEdge(edge); err != nil {
		t.Errorf("Failed to add edge: %v", err)
	}
	e, err := gr.GetEdgeByKey(edge.Key)
	if err != nil || e == nil {
		t.Errorf("Failed to get edge: %v", err)
	}
	if err := gr.RemoveEdgeByKey(edge.Key); err != nil {
		t.Errorf("Failed to remove edge: %v", err)
	}
	if _, err := gr.GetEdgeByKey(edge.Key); err == nil {
		t.Errorf("Expected error getting removed edge")
	}
}

func TestAddEdgeNoDuplicate(t *testing.T) {
	gr := graph.MakeGraph(graph.WithGraphMulti(false))
	src := graph.MakeNode(1)
	dst := graph.MakeNode(2)
	gr.AddNode(src)
	gr.AddNode(dst)
	edge1 := graph.MakeEdge(1, src.Key, dst.Key)
	if err := gr.AddEdge(edge1); err != nil {
		t.Errorf("Failed to add first edge: %v", err)
	}
	edge2 := graph.MakeEdge(2, src.Key, dst.Key)
	if err := gr.AddEdge(edge2); err == nil {
		t.Errorf("Expected error when adding duplicate edge")
	}
}

func TestAddEdgeNodesMustExist(t *testing.T) {
	gr := graph.MakeGraph()
	node := graph.MakeNode(1)
	gr.AddNode(node)
	edge := graph.MakeEdge(1, node.Key, 2)
	if err := gr.AddEdge(edge); err == nil {
		t.Errorf("Expected error when adding edge with non-existing node")
	}
}
