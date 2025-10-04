/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

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
	Nodes   []*Node  `json:"nodes"`
	Edges   []*Edge  `json:"edges"`
	Options TOptions `json:"options"`
}

func MakeGraph(options ...Option[Graph]) *Graph {
	gr := &Graph{}
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

func WithGraphNodes(nodes []*Node) Option[Graph] {
	return func(gr *Graph) {
		gr.Nodes = nodes
	}
}

func WithGraphEdges(edges []*Edge) Option[Graph] {
	return func(gr *Graph) {
		gr.Edges = edges
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
