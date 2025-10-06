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
	"slices"
)

type Option[T any] func(*T) // Type representing functional options pattern

type TKey uint64    // Key type. Can be replaced with any UNIQUE type
type TWeight uint64 // Weight type. Can be replaced with any COMPARABLE type

/*
 * Node struct represents graph node. It has a unique key and optional label.
 * You can properly construct Node via this:
 *
 * node := MakeNode(1)
 *
 * or this:
 *
 * labeledNode := MakeNode(1, WithNodeLabel("Aboba"))
 */

type Node struct {
	Key   TKey   `json:"key"`
	Label string `json:"label"`
}

func MakeNode(key TKey, options ...Option[Node]) *Node {
	node := &Node{}
	node.Key = key
	for _, opt := range options {
		opt(node)
	}
	return node
}

func (node *Node) UpdateNode(options ...Option[Node]) {
	for _, opt := range options {
		opt(node)
	}
}

func WithNodeLabel(label string) Option[Node] {
	return func(node *Node) {
		node.Label = label
	}
}

/*
 * Edge struct represents graph edge -- a connection between 2 nodes.
 * This edge implementation supposed to be direct. If you need to make
 * undirected edge, you should make 2 edges and work with this.
 *
 * Edge represents connection from Edge.Source to Edge.Destination, optionally
 * labelled and weighted.
 *
 * I.e., there are 2 nodes given:
 *
 * src := MakeNode(1)
 * dst := MakeNode(2)
 *
 * You can properly construct this via this:
 *
 * edge := MakeEdge(1, src.Key, dst.Key)
 *
 * or with optional fields:
 *
 * fullyConstructedEdge := MakeEdge(1, src.Key, dst.Key, WithEdgeLabel("Path"), WithEdgeWeight(69))
 */

type Edge struct {
	Key         TKey    `json:"key"`
	Source      TKey    `json:"source"`
	Destination TKey    `json:"destination"`
	Weight      TWeight `json:"weight"`
	Label       string  `json:"label"`
}

func MakeEdge(key, src, dst TKey, options ...Option[Edge]) *Edge {
	edge := &Edge{}
	edge.Key, edge.Source, edge.Destination = key, src, dst
	for _, opt := range options {
		opt(edge)
	}
	return edge
}

func (edge *Edge) UpdateEdge(options ...Option[Edge]) {
	for _, opt := range options {
		opt(edge)
	}
}

func WithEdgeWeight(weight TWeight) Option[Edge] {
	return func(edge *Edge) {
		edge.Weight = weight
	}
}

func WithEdgeLabel(label string) Option[Edge] {
	return func(edge *Edge) {
		edge.Label = label
	}
}

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

func (gr *Graph) UpdateGraph(options ...Option[Graph]) {
	for _, opt := range options {
		opt(gr)
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
		slices.Contains(gr.AdjacencyMap[edge.Source], edge.Destination) ||
		(!gr.Options.IsDirected && slices.Contains(gr.AdjacencyMap[edge.Destination], edge.Source)) {
		return ThrowSameEdgeNotAllowed(edge.Source, edge.Destination)
	}

	if src, _ := gr.GetNodeByKey(edge.Source); src == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Source)
	}

	if dst, _ := gr.GetNodeByKey(edge.Destination); dst == nil {
		return ThrowEdgeEndNotExists(edge.Key, edge.Destination)
	}

	gr.Edges[edge.Key] = edge
	gr.AdjacencyMap[edge.Source] = append(gr.AdjacencyMap[edge.Source], edge.Destination)

	return nil
}

func (gr *Graph) RemoveEdgeByKey(key TKey) error {
	if edge, _ := gr.GetEdgeByKey(key); edge == nil {
		return ThrowEdgeWithKeyNotExists(key)
	}

	delete(gr.Edges, key)
	return nil
}

/*
 * File handling moved to CLI service -- here will be declared just marshalling
 * and unmarshalling handlers
 */

func (g *Graph) MarshalJSON() ([]byte, error) {
	type MarshalGraph Graph
	return json.Marshal(&struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(g),
	})
}

func (g *Graph) UnmarshalJSON(data []byte) error {
	type MarshalGraph Graph
	aux := &struct {
		*MarshalGraph
	}{
		MarshalGraph: (*MarshalGraph)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return ThrowGraphUnmarshalError()
	}
	return nil
}

func (g *Graph) ToJSON() (string, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *Graph) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), g)
}
