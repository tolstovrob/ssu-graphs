/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

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
