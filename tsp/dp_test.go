package solver

import (
	"testing"

	"github.com/nathanchrs/graphalgo"
)

var SimpleSymmetricGraph = graphalgo.AdjListGraph[string, int]{
	"A": {{"B", 20}, {"C", 42}, {"D", 35}},
	"B": {{"A", 20}, {"C", 30}, {"D", 34}},
	"C": {{"A", 42}, {"B", 30}, {"D", 12}},
	"D": {{"A", 35}, {"B", 34}, {"C", 12}},
}

func TestDynamicProgrammingTSPSolver(t *testing.T) {
	testCases := []struct {
		name           string
		graph          graphalgo.AdjListGraph[string, int]
		startVertex    string
		expectedWeight int
		expectedPaths  [][]string
		expectedErr    bool
	}{
		{
			name:           "Simple symmetric graph",
			graph:          SimpleSymmetricGraph,
			startVertex:    "A",
			expectedWeight: 97,
			expectedPaths: [][]string{
				{"A", "B", "C", "D", "A"},
				{"A", "D", "C", "B", "A"},
			},
		},
		{
			name:           "Simple symmetric graph, different starting point",
			graph:          SimpleSymmetricGraph,
			startVertex:    "C",
			expectedWeight: 97,
			expectedPaths: [][]string{
				{"C", "D", "A", "B", "C"},
				{"C", "B", "A", "D", "C"},
			},
		},
	}

	for _, tc := range testCases {
		testGraph := graphalgo.NewReverseAdjListGraphFromAdjListGraph(tc.graph)

		weight, path, err := DynamicProgrammingTSPSolver[string, int](testGraph, tc.startVertex)
		if err != nil && !tc.expectedErr {
			t.Logf("Unexpected error: %v", err)
			t.FailNow()
		} else if err == nil && tc.expectedErr {
			t.Logf("Expecting error but got no error")
			t.FailNow()
		}

		if tc.expectedWeight != weight {
			t.Logf("Weight expected: %d, got: %d", tc.expectedWeight, weight)
			t.Fail()
		}

		hasAnyMatchingPath := false
		for _, expectedPath := range tc.expectedPaths {
			if isPathEqual(expectedPath, path) {
				hasAnyMatchingPath = true
				break
			}
		}
		if !hasAnyMatchingPath {
			t.Logf("Path expected (any): %v, got: %v", tc.expectedPaths, path)
			t.Fail()
		}
	}
}

func isPathEqual[V comparable](a []V, b []V) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
