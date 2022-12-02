# graphalgo

A collection of graph algorithms implemented in Golang.

## Travelling Salesman Problem

Implemented in the `tsp` package.

### Dynamic Programming

Calculates an exact solution. Implemented as the `DynamicProgrammingTSPSolver(graph, startVertex)` function.

- Time complexity: `O(n^2 2^n)`.
- Space complexity: `O(n 2^n)`.
- Can only handle a maximum of 31 vertices.
