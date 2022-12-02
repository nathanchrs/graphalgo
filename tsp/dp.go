package solver

import (
	"errors"
	"fmt"

	graphalgo "github.com/nathanchrs/graphalgo"
)

const MaxInputVertices = 31

type DynamicProgrammingTSPSolverGraph[V comparable, W graphalgo.Weight] interface {
	GetVertices() []V
	GetIncomingEdges(V) []graphalgo.HalfEdge[V, W]
}

func DynamicProgrammingTSPSolver[V comparable, W graphalgo.Weight](graph DynamicProgrammingTSPSolverGraph[V, W], startVertex V) (solutionWeight W, solutionVertices []V, err error) {
	if graph == nil {
		return 0, nil, nil
	}
	vertices := graph.GetVertices()

	// Check vertex count
	if len(vertices) <= 0 {
		return 0, nil, nil
	}
	if len(vertices) == 1 {
		return 0, []V{vertices[0]}, nil
	}
	if len(vertices) > MaxInputVertices {
		return 0, nil, fmt.Errorf("DynamicProgrammingTSPSolver can only handle input graphs with up to %d vertices", MaxInputVertices)
	}

	// Prepare mapping between vertex and vertex index
	vertexIdxMap := make(map[V]int, len(vertices))
	for i, v := range vertices {
		vertexIdxMap[v] = i
	}
	if len(vertexIdxMap) != len(vertices) {
		return 0, nil, fmt.Errorf("DynamicProgrammingTSPSolver: graph.GetVertices() must return unique vertices")
	}
	startVertexIdx, ok := vertexIdxMap[startVertex]
	if !ok {
		return 0, nil, fmt.Errorf("DynamicProgrammingTSPSolver: unknown startVertex %v", startVertex)
	}
	getIncomingEdges := func(toIdx int) ([]incomingEdge[W], error) {
		to := vertices[toIdx]
		incomingEdges := graph.GetIncomingEdges(to)
		incomingEdgesWithIdx := make([]incomingEdge[W], len(incomingEdges))
		for i, e := range incomingEdges {
			fromIdx, ok := vertexIdxMap[e.Vertex]
			if !ok {
				return nil, fmt.Errorf("DynamicProgrammingTSPSolver: unknown vertex %v returned by graph.GetIncomingEdges(%v)", e.Vertex, to)
			}
			incomingEdgesWithIdx[i] = incomingEdge[W]{
				FromIdx: fromIdx,
				Weight:  e.Weight,
			}
		}
		return incomingEdgesWithIdx, nil
	}

	memo := make([][]*W, CountBitArrangements(len(vertices)))

	// Run the dynamic programming solver to find the smallest total weight
	solutionWeight, err = tspDP(memo, len(vertices), getIncomingEdges, startVertexIdx, NewBitsetAllSet(len(vertices)), startVertexIdx)
	if err != nil {
		if !errors.Is(err, graphalgo.ErrNoSolution) {
			return 0, nil, err
		}
	}

	// Backtrack memo to find the vertices that is part of the optimal solution
	solutionVertexIdx, err := tspDPPath(memo, len(vertices), getIncomingEdges, startVertexIdx, NewBitsetAllSet(len(vertices)), startVertexIdx)
	if err != nil {
		return 0, nil, err
	}
	for _, vertexIdx := range solutionVertexIdx {
		solutionVertices = append(solutionVertices, vertices[vertexIdx])
	}

	return solutionWeight, solutionVertices, nil
}

type incomingEdge[W graphalgo.Weight] struct {
	FromIdx int
	Weight  W
}

func tspDP[W graphalgo.Weight](memo [][]*W, vertexCount int, getIncomingEdges func(toIdx int) ([]incomingEdge[W], error), startVertexIdx int, visited Bitset, lastVertexIdx int) (W, error) {
	// Base case
	if !visited.IsAnySet() {
		if lastVertexIdx != startVertexIdx {
			panic("DynamicProgrammingTSPSolver: base case must start at the starting vertex")
		}
		return 0, nil
	}

	// Use memoized result if it has been calculated before
	if memo[visited] != nil && memo[visited][lastVertexIdx] != nil {
		return *memo[visited][lastVertexIdx], nil
	}

	var result *W

	incomingEdges, err := getIncomingEdges(lastVertexIdx)
	if err != nil {
		return 0, err
	}

	for _, incomingEdge := range incomingEdges {
		visitedWithoutLastVertex := visited.Unset(lastVertexIdx)

		if visited.IsSet(incomingEdge.FromIdx) || (!visitedWithoutLastVertex.IsAnySet() && incomingEdge.FromIdx == startVertexIdx) {
			weight, err := tspDP(memo, vertexCount, getIncomingEdges, startVertexIdx, visitedWithoutLastVertex, incomingEdge.FromIdx)
			if err != nil {
				if !errors.Is(err, graphalgo.ErrNoSolution) {
					return 0, err
				}
				continue
			}

			weight += incomingEdge.Weight
			if result == nil || weight < *result {
				result = &weight
			}
		}
	}

	// Return error if no solution found
	if result == nil {
		return 0, graphalgo.ErrNoSolution
	}

	// Save result in memo
	if memo[visited] == nil {
		memo[visited] = make([]*W, vertexCount)
	}
	memo[visited][lastVertexIdx] = result

	return *result, nil
}

func tspDPPath[W graphalgo.Weight](memo [][]*W, vertexCount int, getIncomingEdges func(toIdx int) ([]incomingEdge[W], error), startVertexIdx int, visited Bitset, lastVertexIdx int) ([]int, error) {
	// Base case
	if !visited.IsAnySet() {
		if lastVertexIdx != startVertexIdx {
			panic("DynamicProgrammingTSPSolver: base case must start at the starting vertex")
		}
		return []int{startVertexIdx}, nil
	}

	var prevVisited Bitset
	var prevLastVertexIdx int
	var minWeight *W

	incomingEdges, err := getIncomingEdges(lastVertexIdx)
	if err != nil {
		return nil, err
	}

	// Find previous vertex
	for _, incomingEdge := range incomingEdges {
		visitedWithoutLastVertex := visited.Unset(lastVertexIdx)

		if visited.IsSet(incomingEdge.FromIdx) || (!visitedWithoutLastVertex.IsAnySet() && incomingEdge.FromIdx == startVertexIdx) {
			var prevWeight *W
			if memo[visitedWithoutLastVertex] != nil {
				prevWeight = memo[visitedWithoutLastVertex][incomingEdge.FromIdx]
			} else if !visitedWithoutLastVertex.IsAnySet() {
				var zeroW W
				prevWeight = &zeroW
			}
			if prevWeight == nil {
				panic("DynamicProgrammingTSPSolver: invalid backtrack, possibly caused by inconsistent input graph")
			}
			weight := *prevWeight + incomingEdge.Weight
			if minWeight == nil || weight < *minWeight {
				minWeight = &weight
				prevVisited = visitedWithoutLastVertex
				prevLastVertexIdx = incomingEdge.FromIdx
			}
		}
	}

	// Backtrack to previous vertex
	path, err := tspDPPath(memo, vertexCount, getIncomingEdges, startVertexIdx, prevVisited, prevLastVertexIdx)
	if err != nil {
		return nil, err
	}

	return append(path, lastVertexIdx), nil
}
