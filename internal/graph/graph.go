package graph

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type WeightedGraph struct {
	nodes map[string]map[string]int
	lock  sync.RWMutex
}

type Path struct {
	steps  []string
	length int
}

func NewWeightedGraph() *WeightedGraph {
	return &WeightedGraph{nodes: make(map[string]map[string]int)}
}

func NewPath() *Path {
	return &Path{}
}

func (g *WeightedGraph) GetAllPathsToDepth(depth int) []*Path {
	var paths []*Path
	var frontier []*Path
	visited := make(map[string]bool)
	// Initialize the frontier with all 1 word paths
	g.lock.RLock()
	for k := range g.nodes {
		path := NewPath()
		path.AppendStep(k, 1)
		frontier = append(frontier, path)
	}
	g.lock.RUnlock()

	fmt.Println(len(g.nodes))
	for len(frontier) > 0 {
		exploredPath, newFrontier := frontier[0], frontier[1:]
		visited[exploredPath.PathHash()] = true
		if exploredPath.NumSteps() <= depth {
			paths = append(paths, exploredPath)
		}
		previousStep, err := exploredPath.LastStep()
		if err != nil { // TODO: Get better at error handling
			panic("Frig")
		}
		neighbors, err := g.GetNeighbors(previousStep)
		if err != nil {
			panic("Frig")
		}

		g.lock.RLock()
		for word := range neighbors {
			weight, ok := neighbors[word]
			if !ok {
				panic("Invalid node")
			}
			newPath := AddStep(exploredPath, word, weight)
			_, ok = visited[newPath.PathHash()]
			if !ok {
				newFrontier = append(newFrontier, newPath)
			}
		}
		g.lock.RUnlock()
		fmt.Println("Done with search")
		frontier = newFrontier
	}
	return paths
}

func (g *WeightedGraph) AddEdge(start string, end string) {
	g.addNode(end)
	g.addNode(start)
	g.updateNodePaths(start, end)
}

func (g *WeightedGraph) GetNeighbors(node string) (map[string]int, error) {
	g.lock.RLock()
	neigbors, ok := g.nodes[node]
	g.lock.RUnlock()
	if !ok {
		return nil, errors.New("Node not found")
	}
	return neigbors, nil
}

func (g *WeightedGraph) addNode(node string) {
	g.lock.Lock()
	_, ok := g.nodes[node]
	if !ok {
		g.nodes[node] = make(map[string]int)
	}
	g.lock.Unlock()
}

func (g *WeightedGraph) updateNodePaths(start string, end string) {
	g.lock.RLock()
	paths, ok := g.nodes[start]
	g.lock.RUnlock()

	g.lock.Lock()

	if ok {
		weight, ok := paths[end]
		if ok {
			paths[end] = weight + 1
		} else {
			paths[end] = 1
		}
	}
	g.lock.Unlock()
}

func (path *Path) AppendStep(word string, distance int) {
	path.steps = append(path.steps, word)
	path.length = path.length + distance
}

func AddStep(path *Path, word string, distance int) *Path {
	return &Path{
		steps:  append(path.steps, word),
		length: path.length + distance,
	}
}

func (path *Path) NumSteps() int {
	return len(path.steps)
}

func (path *Path) Distance() int {
	return path.length
}

func (path *Path) LastStep() (string, error) {
	if path.NumSteps() == 0 {
		return "", errors.New("Empty Path")
	} else {
		return path.steps[len(path.steps)-1], nil
	}
}

func (path *Path) PrintPath() {
	for _, word := range path.steps {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
	fmt.Println(path.length)
}

func (path *Path) PathHash() string {
	return strings.Join(path.steps, ", ")
}

func (path *Path) Steps() []string {
	return path.steps
}

func (path *Path) Length() int {
	return path.length
}
