/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

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
