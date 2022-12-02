package graphalgo

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var ErrNoSolution = errors.New("no solution")

type Weight interface {
	constraints.Integer | constraints.Float
}

type Edge[V comparable, W Weight] struct {
	From   V
	To     V
	Weight W
}

type HalfEdge[V comparable, W Weight] struct {
	Vertex V
	Weight W
}

type AdjListGraph[V comparable, W Weight] map[V][]HalfEdge[V, W]

func NewAdjListGraphFromReverseAdjListGraph[V comparable, W Weight](g ReverseAdjListGraph[V, W]) AdjListGraph[V, W] {
	res := make(AdjListGraph[V, W])
	for v, edges := range g {
		if len(edges) == 0 {
			res[v] = nil // Preserve vertices with no adjacent edges
		}
		for _, e := range edges {
			resEdges := res[e.Vertex]
			resEdges = append(resEdges, HalfEdge[V, W]{Vertex: v, Weight: e.Weight})
			res[e.Vertex] = resEdges
		}
	}
	return res
}

func (g AdjListGraph[V, W]) GetVertices() []V {
	if g == nil {
		return nil
	}
	var vertices []V
	for v, _ := range g {
		vertices = append(vertices, v)
	}
	return vertices
}

func (g AdjListGraph[V, W]) GetOutgoingEdges(from V) []HalfEdge[V, W] {
	if g == nil {
		return nil
	}
	outgoingEdges, ok := g[from]
	if !ok {
		return nil
	}
	return outgoingEdges
}

type ReverseAdjListGraph[V comparable, W Weight] AdjListGraph[V, W]

func NewReverseAdjListGraphFromAdjListGraph[V comparable, W Weight](g AdjListGraph[V, W]) ReverseAdjListGraph[V, W] {
	input := ReverseAdjListGraph[V, W](g)
	output := NewAdjListGraphFromReverseAdjListGraph[V, W](input)
	return ReverseAdjListGraph[V, W](output)
}

func (g ReverseAdjListGraph[V, W]) GetVertices() []V {
	return AdjListGraph[V, W](g).GetVertices()
}

func (g ReverseAdjListGraph[V, W]) GetIncomingEdges(to V) []HalfEdge[V, W] {
	return AdjListGraph[V, W](g).GetOutgoingEdges(to)
}
