package main

import (
	"fmt"
	. "hive-arena/common"
)

// BuildWallOrders: deterministic wall-building strategy
// - For each non-carrying bee (in order) try to find a visible enemy hive
// - For each hive try directions; compute mid (tile between hive and outward)
//   and anchor (tile where bee must stand = mid.Neighbour(dir))
// - Require both mid and anchor to be known, walkable and empty
// - Pathfind bee to anchor; if on anchor and resources available, BUILD_WALL toward mid
// - Otherwise MOVE one step toward anchor
// - Reserve mid tiles so only a single bee aims to create a wall there this turn
// - Stop after issuing a single BUILD_WALL in this turn (to reduce conflicts)
func BuildWallOrders(state *GameState, player int, as *AgentState) []Order {
	orders := []Order{}

	// reservation map for mid tiles (to avoid multiple bees targeting same wall)
	reserved := make(map[Coords]bool)
	builtThisTurn := false

	for _, bee := range as.MyBees {
		if bee.HasFlower {
			continue
		}

		if builtThisTurn {
			break
		}

		// iterate enemy hives we remember
		for owner, hiveList := range as.Hives {
			if owner == player {
				continue
			}
			for _, hive := range hiveList {
				if builtThisTurn {
					break
				}
				// try every direction outward from hive
				for dir := range DirectionToOffset {
					// mid is tile adjacent to hive; wall will be placed here
					mid := hive.Neighbour(dir)
					// anchor is the tile the bee must stand on to build into mid
					anchor := mid.Neighbour(dir)

					// skip if we've already reserved or built here
					if reserved[mid] || builtThisTurn {
						continue
					}

					// both tiles must exist in the authoritative state
					midHex, ok1 := state.Hexes[mid]
					anchorHex, ok2 := state.Hexes[anchor]
					if !ok1 || !ok2 {
						continue
					}

					// must be walkable
					if !midHex.Terrain.IsWalkable() || !anchorHex.Terrain.IsWalkable() {
						continue
					}

					// mid must be empty (no entity) at planning time
					if midHex.Entity != nil {
						continue
					}

					// anchor is allowed to be occupied only if it's this bee itself
					if anchorHex.Entity != nil && anchorHex.Entity.Player != bee.Player {
						continue
					}

					// skip if there's already a recorded wall in memory
					if as.IsWall(mid) {
						continue
					}

					// pathfind bee to anchor
					path, ok := as.find_path(bee, anchor)
					if !ok || len(path) == 0 {
						continue
					}

					// if bee is already on anchor, attempt build (if resources available)
					if bee.Coords == anchor {
						if int(state.PlayerResources[player]) >= int(WALL_COST) {
							// compute direction from bee->mid
							var dirToMid Direction
							found := false
							for d := range DirectionToOffset {
								if bee.Coords.Neighbour(d) == mid {
									dirToMid = d
									found = true
									break
								}
							}
							if !found {
								continue
							}

							// final safety check: mid still empty in state
							if state.Hexes[mid] != nil && state.Hexes[mid].Entity == nil {
								orders = append(orders, Order{Type: BUILD_WALL, Coords: bee.Coords, Direction: dirToMid})
								reserved[mid] = true
								builtThisTurn = true
								fmt.Printf("WallBuilder: BUILD_WALL planned by bee %v at anchor %v into mid %v (hive %v)\n", bee.Coords, anchor, mid, hive)
								break
							}
						}
						// else (insufficient resources) skip
						continue
					}

					// otherwise issue single-step MOVE along path toward anchor
					if len(path) > 1 {
						next := path[1]
						if dirMove, ok := as.BestDirectionTowards(bee.Coords, next); ok {
							if _, ok2 := IsValidMoveTarget(as, bee.Coords, dirMove); ok2 {
								orders = append(orders, Order{Type: MOVE, Coords: bee.Coords, Direction: dirMove})
								fmt.Printf("WallBuilder: MOVE planned for bee %v toward anchor %v (next=%v)\n", bee.Coords, anchor, next)
								// reserve mid to avoid other bees targeting same wall location this turn
								reserved[mid] = true
								break
							}
						}
					}
				}
			}
		}
	}

	return orders
}