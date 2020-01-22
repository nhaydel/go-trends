package graph

import (
	"testing"
)

func TestAddEdge(test *testing.T) {
	graph := NewWeightedGraph()
	graph.AddEdge("first", "second")
	graph.AddEdge("first", "second")

	first_neighbors, _ := graph.GetNeighbors("first")
	second_neighbors, _ := graph.GetNeighbors("second")

	for neighbor, weight := range first_neighbors {
		if weight != 2 {
			test.Errorf("Incorrect weight for %s. Expected 2, got %d", neighbor, weight)
		}
	}

	if len(second_neighbors) != 0 {
		test.Errorf("Node %s has an incorrect number of neighbors. Expected 0, got %d", "second", len(second_neighbors))
	}
}

func TestGetAllPaths(test *testing.T) {
	graph := NewWeightedGraph()
	graph.AddEdge("first", "second")
	graph.AddEdge("first", "second")
	paths := graph.GetAllPathsToDepth(5)
	for _, path := range paths {
		if path.NumSteps() > 5 {
			test.Errorf("Path was greater than the specified length.")
		}
	}

}
