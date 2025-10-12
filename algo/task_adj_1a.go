/*
 * This package contains algorithms and tasks for my SSU course
 */

package algo

import "github.com/tolstovrob/graph-go/graph"

/*
 * Task: Get all nodes, for which degree is greater then half-degree of entry
 */

func InDegreeLessThan(gr *graph.Graph, targetKey graph.TKey) []graph.TKey {
	// Find all entries
	inDegree := make(map[graph.TKey]int)
	for _, edge := range gr.Edges {
		inDegree[edge.Destination]++
		if !gr.Options.IsDirected {
			inDegree[edge.Source]++
		}
	}
	targetInDegree := inDegree[targetKey]

	// Find all nodes satisfies task objective
	var result []graph.TKey
	for nodeKey := range gr.Nodes {
		if inDegree[nodeKey] < targetInDegree {
			result = append(result, nodeKey)
		}
	}

	return result
}
