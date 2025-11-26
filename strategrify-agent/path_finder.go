package main

// should get the known map & goal where to go
// returns an array or list of shortest path to the goal
// or if not path found then nil?
func (as *AgentState) BFS(start, goal Coords) ([]Coords, bool) {
    type node struct {
        Hex_c	Coords
        Prev	*node
    }

    visited := make(map[Coords]bool) //(x,y), known?

    // Use a queue of *node so Prev pointers are stable
    startNode := &node{Hex_c: start, Prev: nil}
    queue := []*node{startNode} // slice array that holds pointers of nodes (currently we just add one startNode there)
    visited[start] = true      // finds the start node from the map and marks it as true, so we have visited there

    for len(queue) > 0 // as long as we have terrain/map to go through in the queue
	{
        current := queue[0] // we take the first from the queue
        queue = queue[1:] // slices the queue from the second to everything else of the array effective popping the first off from array

        if current.Hex_c == goal { // if the current terrain is the goal
            // reconstruct path by following Prev pointers
            var path []Coords // this is slice not array
            for n := current; n != nil; n = n.Prev {
                path = append([]Coords{n.Hex_c}, path...) // append to the slice always the current node, and always the previous to the front of the old effectively making it from start -> end
            }
            return path, true
        }

        for dir := range DirectionToOffset {
            next := current.Hex_c.Neighbour(dir)

            // Only allow path on known tiles:
            terrain, ok := as.Map[next]
            if !ok {
                continue // unknown tile -> cannot use (change if you want to allow exploring unknowns)
            }
            if terrain == ROCK {
                continue // not walkable
            }
			// can we stand on hive?
			// we also need to check if there is a wall on a terrain and go around it

            if !visited[next] { // if next isnt in visited array?
                visited[next] = true // ad it there and make it true so we know it has been visited?
                nextNode := &node{Hex_c: next, Prev: current}
                queue = append(queue, nextNode) // should this be queue ... ?
            }
        }
    }

    return nil, false // no path found
}