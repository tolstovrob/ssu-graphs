/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

import (
	"encoding/json"
	"fmt"
	"slices"
)

/*
 * Graph struct.
 *
 * First of all, we got TOptions struct, which represents all possible graph
 * configuration. Now, it only has IsMulti and IsDirected for multigraphs and
 * Directed graphs respectively, but it is easily scalable for other options
 * if neccessary.
 *
 * Graph struct sa it is contains Nodes and Edges lists of Node and Edge
 * pointers respectively, and Options configuration of TOptions.
 *
 * Graph represented via adjacency list of edges. But it also possible to have
 * islands with no connections. You cannot find them in Edges, but in the Nodes.
 *
 * You actually can use default constructor with this one. It will build
 * non-multi undirected graph:
 *
 * gr := Graph{}
 *
 * But to make graph properly, use constructor with options:
 *
 * gr := MakeGraph(WithGraphMulti(true), WithGraphDirected(false))
 *
 * I.e., code above will create undirected multigraph.
 */

type TOptions struct {
	IsMulti    bool `json:"isMulti"`
	IsDirected bool `json:"IsDirected"`
}

type Graph struct {
	Nodes        map[TKey]*Node  `json:"nodes"`
	Edges        map[TKey]*Edge  `json:"edges"`
	AdjacencyMap map[TKey][]TKey `json:"adjacencyMap"`
	Options      TOptions        `json:"options"`
}

func MakeGraph(options ...Option[Graph]) *Graph {
	gr := &Graph{}
	gr.Nodes = make(map[TKey]*Node)
	gr.Edges = make(map[TKey]*Edge)
	gr.AdjacencyMap = make(map[TKey][]TKey)
	for _, opt := range options {
		opt(gr)
	}
	return gr
}

func (gr *Graph) RebuildEdges() {
	newEdges := make(map[TKey]*Edge)
	edgeKeysUsed := make(map[TKey]bool)
	edgeKeyCounter := TKey(1)

	nextEdgeKey := func() TKey {
		for {
			if !edgeKeysUsed[edgeKeyCounter] {
				edgeKeysUsed[edgeKeyCounter] = true
				key := edgeKeyCounter
				edgeKeyCounter++
				return key
			}
			edgeKeyCounter++
		}
	}

	edgeID := func(src, dst TKey) string {
		if !gr.Options.IsDirected && src > dst {
			src, dst = dst, src
		}
		return fmt.Sprintf("%d-%d", src, dst)
	}

	seenEdges := make(map[string]bool)

	for _, edge := range gr.Edges {
		id := edgeID(edge.Source, edge.Destination)
		if !gr.Options.IsMulti {
			if seenEdges[id] {
				continue
			}
			seenEdges[id] = true
		}

		key := edge.Key
		if key == 0 || edgeKeysUsed[key] {
			key = nextEdgeKey()
		} else {
			edgeKeysUsed[key] = true
		}

		newEdge := &Edge{
			Key:         key,
			Source:      edge.Source,
			Destination: edge.Destination,
			Weight:      edge.Weight,
			Label:       edge.Label,
		}
		newEdges[key] = newEdge
	}

	gr.Edges = newEdges
}

func (gr *Graph) RebuildAdjacencyMap() {
	gr.AdjacencyMap = make(map[TKey][]TKey)
	for _, edge := range gr.Edges {
		gr.AdjacencyMap[edge.Source] = append(gr.AdjacencyMap[edge.Source], edge.Destination)
		if !gr.Options.IsDirected {
			gr.AdjacencyMap[edge.Destination] = append(gr.AdjacencyMap[edge.Destination], edge.Source)
		}
	}
}

/*
 * Later in the code, Graph.RebuildAdjacencyMap will be called many times.
 * It could really affect performance on huge amount of edges, but since
 * it is just academical example, we will pretend it never happens.
 *
 * Anyways, need to fix, so mark this part as WIP!
 */

func (gr *Graph) UpdateGraph(options ...Option[Graph]) {
	oldOptions := gr.Options

	for _, opt := range options {
		opt(gr)
	}

	if oldOptions != gr.Options {
		gr.RebuildEdges()
		gr.RebuildAdjacencyMap()
	}
}

func WithGraphNodes(nodes map[TKey]*Node) Option[Graph] {
	return func(gr *Graph) {
		gr.Nodes = nodes
	}
}

func WithGraphEdges(edges map[TKey]*Edge) Option[Graph] {
	return func(gr *Graph) {
		gr.Edges = edges
	}
}
func WithGraphAdjacencyMap(adj map[TKey][]TKey) Option[Graph] {
	return func(gr *Graph) {
		gr.AdjacencyMap = adj
	}
}

func WithGraphOptions(options TOptions) Option[Graph] {
	return func(gr *Graph) {
		gr.Options = options
	}
}

func WithGraphMulti(isMulti bool) Option[Graph] {
	return func(gr *Graph) {
		gr.Options.IsMulti = isMulti
	}
}

func WithGraphDirected(IsDirected bool) Option[Graph] {
	return func(gr *Graph) {
		gr.Options.IsDirected = IsDirected
	}
}

/*
 * Next coming finding, adding and removing handlers for nodes and edges. I put
 * them apart main Graph struct because they contain both node and edges and
 * graph. All of them will throw an error if the operation is not allowed (I.e.
 * adding existing node or connecting nodes with more then one time in multi).
 */

func (gr *Graph) GetNodeByKey(key TKey) (*Node, error) {
	if gr.Nodes == nil {
		return nil, ThrowNodesListIsNil()
	}

	if _, exists := gr.Nodes[key]; !exists {
		return nil, ThrowNodeWithKeyNotExists(key)
	}

	return gr.Nodes[key], nil
}

func (gr *Graph) AddNode(node *Node) error {
	if node, _ := gr.GetNodeByKey(node.Key); node != nil {
		return ThrowNodeWithKeyExists(node.Key)
	}

	gr.Nodes[node.Key] = node
	return nil
}

func (gr *Graph) RemoveNodeByKey(key TKey) error {
	if _, err := gr.GetNodeByKey(key); err != nil {
		return err
	}

	delete(gr.Nodes, key)

	for _, edge := range gr.Edges {
		if edge.Source == key || edge.Destination == key {
			gr.RemoveEdgeByKey(edge.Key)
		}
	}

	gr.RebuildAdjacencyMap()
	return nil
}

func (gr *Graph) GetEdgeByKey(key TKey) (*Edge, error) {
	if gr.Edges == nil {
		return nil, ThrowEdgesListIsNil()
	}

	if _, exists := gr.Edges[key]; !exists {
		return nil, ThrowEdgeWithKeyNotExists(key)
	}

	return gr.Edges[key], nil
}

func (gr *Graph) AddEdge(edge *Edge) error {
	if edge, _ := gr.GetEdgeByKey(edge.Key); edge != nil {
		return ThrowEdgeWithKeyExists(edge.Key)
	}

	if !gr.Options.IsMulti &&
		(slices.Contains(gr.AdjacencyMap[edge.Source], edge.Destination) ||
			!gr.Options.IsDirected && slices.Contains(gr.AdjacencyMap[edge.Destination], edge.Source)) {
		return ThrowSameEdgeNotAllowed(edge.Source, edge.Destination)
	}

	if src, _ := gr.GetNodeByKey(edge.Source); src == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Source)
	}

	if dst, _ := gr.GetNodeByKey(edge.Destination); dst == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Destination)
	}

	gr.Edges[edge.Key] = edge
	gr.RebuildAdjacencyMap()
	return nil
}

func (gr *Graph) RemoveEdgeByKey(key TKey) error {
	if edge, _ := gr.GetEdgeByKey(key); edge == nil {
		return ThrowEdgeWithKeyNotExists(key)
	}

	delete(gr.Edges, key)
	gr.RebuildAdjacencyMap()
	return nil
}

/*
 * File handling moved to CLI service -- here will be declared just marshalling
 * and unmarshalling handlers
 */

func (gr *Graph) MarshalJSON() ([]byte, error) {
	type MarshalGraph Graph
	return json.Marshal(&struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(gr),
	})
}

func (gr *Graph) UnmarshalJSON(data []byte) error {
	type MarshalGraph Graph
	aux := &struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(gr),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return ThrowGraphUnmarshalError()
	}
	gr.RebuildAdjacencyMap()
	return nil
}

func (gr *Graph) ToJSON() (string, error) {
	b, err := json.Marshal(gr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (gr *Graph) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), gr)
}
