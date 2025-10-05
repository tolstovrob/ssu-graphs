package graph

import "fmt"

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
