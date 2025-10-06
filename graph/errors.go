package graph

import "fmt"

func ThrowNodesListIsNil() error {
	return fmt.Errorf("Nodes list is nil")
}

func ThrowEdgesListIsNil() error {
	return fmt.Errorf("Edges list is nil")
}

func ThrowNodeWithKeyExists(key TKey) error {
	return fmt.Errorf("Node with key %v already exists", key)
}

func ThrowNodeWithKeyNotExists(key TKey) error {
	return fmt.Errorf("Node with key %v not exists", key)
}

func ThrowEdgeWithKeyExists(key TKey) error {
	return fmt.Errorf("Edge with key %v already exists", key)
}

func ThrowEdgeWithKeyNotExists(key TKey) error {
	return fmt.Errorf("Edge with key %v not exists", key)
}

func ThrowSameEdgeNotAllowed(src, dst TKey) error {
	return fmt.Errorf("Edge with src: %v and dst: %v already exists. If you don't think so, check your graph's options", src, dst)
}

func ThrowEdgeEndNotExists(key TKey, end TKey) error {
	return fmt.Errorf("Edge %v has end %v, which is not represented in Nodes", key, end)
}

func ThrowGraphUnmarshalError() error {
	return fmt.Errorf("Cannot unmarshal graph. Seems like this is a internal error")
}
