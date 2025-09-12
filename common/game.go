package common

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

const (
	INIT_FIELD_FLOWERS = 60
	INIT_HIVE_HP       = 12
	INIT_BEE_HP        = 2
	INIT_WALL_HP       = 6

	BEE_COST  = 12
	HIVE_COST = 24
	WALL_COST = 6

	HIVE_FIELD_OF_VIEW = 4
	INFLUENCE_TIMEOUT  = 50
)

type Entity struct {
	Type   EntityType	`json:"type"`
	HP     int			`json:"hp"`
	Player int		`json:"player"`
}

type EntityType int

const (
	WALL EntityType = iota
	HIVE
	BEE
)

type Hex struct {
	Terrain   Terrain	`json:"terrain"`
	Resources uint		`json:"resources,omitempty"`
	Influence int	`json:"influence"`
	Entity    *Entity	`json:"entity,omitempty"`
}

type Order struct {
	Type      OrderType
	Player    int
	Coords    Coords
	Direction Direction

	Status OrderStatus
}

type OrderType int

const (
	MOVE OrderType = iota
	ATTACK
	BUILD_WALL
	BUILD_HIVE
	FORAGE
	SPAWN
)

type OrderStatus int

const (
	PENDING OrderStatus = iota
	INVALID_UNIT
	BLOCKED
	INVALID_TARGET
	CANNOT_FORAGE
	NOT_ENOUGH_RESOURCES
	UNIT_ALREADY_ACTED
	OK
)

func (o *Order) UnitType() EntityType {
	if o.Type == SPAWN {
		return HIVE
	}
	return BEE
}

func (o *Order) Target() Coords {
	return o.Coords.Neighbour(o.Direction)
}

type GameState struct {
	NumPlayers          int		`json:"numPlayers"`
	Turn                uint			`json:"turn"`
	Hexes               map[Coords]*Hex	`json:"hexes"`
	PlayerResources     []uint			`json:"playerResources"`
	LastInfluenceChange uint			`json:"lastInfluenceChange"`

	Winners  []int					`json:"winners,omitempty"`
	GameOver bool						`json:"gameOver"`
}

var playerMappings = [][]int{
	{},
	{0, -1, -1, -1, -1, -1},
	{0, -1, -1, 1, -1, -1},
	{0, -1, 1, -1, 2, -1},
	{-1, 0, 1, -1, 2, 3},
	{0, 1, 2, 3, 4, -1},
	{0, 1, 2, 3, 4, 5},
}

func NewGameState(mapData *MapData, numPlayers int) *GameState {
	gs := &GameState{
		NumPlayers: numPlayers,
		Hexes:      make(map[Coords]*Hex),
	}

	for coords, terrain := range mapData.Map {
		gs.Hexes[coords] = &Hex{Terrain: terrain}
	}

	for _, spawn := range mapData.Spawns {
		player := playerMappings[numPlayers][spawn.Player]
		if player == -1 {
			continue
		}

		switch spawn.Kind {
		case HIVE:
			gs.Hexes[spawn.Coords].Entity = &Entity{Type: HIVE, HP: INIT_HIVE_HP, Player: player}
		case BEE:
			gs.Hexes[spawn.Coords].Entity = &Entity{Type: BEE, HP: INIT_BEE_HP, Player: player}
		}
	}

	for _, hex := range gs.Hexes {
		if hex.Terrain == FIELD {
			hex.Resources = INIT_FIELD_FLOWERS
		}
	}

	gs.PlayerResources = make([]uint, numPlayers)
	gs.updateInfluence()
	gs.checkEndGame()

	return gs
}

func (gs *GameState) EntityAt(coords Coords) *Entity {
	hex, ok := gs.Hexes[coords]
	if !ok {
		return nil
	}
	return hex.Entity
}

func (gs *GameState) TerrainAt(coords Coords) Terrain {
	hex, ok := gs.Hexes[coords]
	if !ok {
		return INVALID
	}
	return hex.Terrain
}

func (gs *GameState) Entities() []struct {
	Coords Coords
	Entity *Entity
} {
	var entities []struct {
		Coords Coords
		Entity *Entity
	}
	for coords, hex := range gs.Hexes {
		if hex.Entity != nil {
			entities = append(entities, struct {
				Coords Coords
				Entity *Entity
			}{coords, hex.Entity})
		}
	}
	return entities
}

// ProcessOrders processes a set of orders for all players.
func (gs *GameState) ProcessOrders(orders [][]Order) ([]Order, error) {
	if gs.GameOver {
		return nil, fmt.Errorf("cannot process orders in a finished game")
	}

	for id := int(0); id < gs.NumPlayers; id++ {
		for i := range orders[id] {
			orders[id][i].Player = id
		}
	}

	var rounds [][]Order
	if len(orders) > 0 {
		maxLen := 0
		for _, playerOrders := range orders {
			if len(playerOrders) > maxLen {
				maxLen = len(playerOrders)
			}
		}

		rounds = make([][]Order, maxLen)
		for i := 0; i < maxLen; i++ {
			for j := range orders {
				if i < len(orders[j]) {
					rounds[i] = append(rounds[i], orders[j][i])
				}
			}
		}
	}

	acted := make(map[*Entity]bool)
	var processed []Order

	for _, round := range rounds {
		rand.Shuffle(len(round), func(i, j int) {
			round[i], round[j] = round[j], round[i]
		})

		for i := range round {
			order := &round[i]
			unit := gs.EntityAt(order.Coords)
			if unit != nil {
				if acted[unit] {
					order.Status = UNIT_ALREADY_ACTED
					continue
				}
				gs.applyOrder(order)
				acted[unit] = true
				processed = append(processed, *order)
			}
		}
	}

	gs.Turn++
	gs.updateInfluence()
	gs.checkEndGame()

	return processed, nil
}

func (gs *GameState) applyOrder(order *Order) {
	switch order.Type {
	case MOVE:
		gs.applyMoveOrder(order)
	case ATTACK:
		gs.applyAttackOrder(order)
	case BUILD_WALL:
		gs.applyBuildWallOrder(order)
	case BUILD_HIVE:
		gs.applyBuildHiveOrder(order)
	case FORAGE:
		gs.applyForageOrder(order)
	case SPAWN:
		gs.applySpawnOrder(order)
	}
}

func (gs *GameState) getUnit(order *Order) *Entity {
	unit := gs.EntityAt(order.Coords)
	if unit == nil || unit.Type != order.UnitType() || unit.Player != order.Player {
		order.Status = INVALID_UNIT
		return nil
	}
	return unit
}

func (gs *GameState) targetIsBlocked(order *Order) bool {
	targetTerrain := gs.TerrainAt(order.Target())
	entity := gs.EntityAt(order.Target())

	if !targetTerrain.IsWalkable() || entity != nil {
		order.Status = BLOCKED
		return true
	}
	return false
}

func (gs *GameState) tryToPay(order *Order, cost uint) bool {
	if gs.PlayerResources[order.Player] < cost {
		order.Status = NOT_ENOUGH_RESOURCES
		return false
	}
	gs.PlayerResources[order.Player] -= cost
	return true
}

func (gs *GameState) applyMoveOrder(order *Order) {
	bee := gs.getUnit(order)
	if bee == nil {
		return
	}
	if gs.targetIsBlocked(order) {
		return
	}

	gs.Hexes[order.Coords].Entity = nil
	gs.Hexes[order.Target()].Entity = bee

	order.Status = OK
}

func (gs *GameState) applyAttackOrder(order *Order) {
	if gs.getUnit(order) == nil {
		return
	}

	entity := gs.EntityAt(order.Target())
	if entity == nil {
		order.Status = INVALID_TARGET
		return
	}

	entity.HP--
	if entity.HP <= 0 {
		gs.Hexes[order.Target()].Entity = nil
	}

	order.Status = OK
}

func (gs *GameState) applyBuildWallOrder(order *Order) {
	if gs.getUnit(order) == nil {
		return
	}
	if gs.targetIsBlocked(order) {
		return
	}
	if !gs.tryToPay(order, WALL_COST) {
		return
	}

	wall := &Entity{Type: WALL, HP: INIT_WALL_HP, Player: order.Player}
	gs.Hexes[order.Target()].Entity = wall

	order.Status = OK
}

func (gs *GameState) applyBuildHiveOrder(order *Order) {
	if gs.getUnit(order) == nil {
		return
	}
	if !gs.tryToPay(order, HIVE_COST) {
		return
	}

	hive := &Entity{Type: HIVE, HP: INIT_HIVE_HP, Player: order.Player}
	gs.Hexes[order.Coords].Entity = hive

	order.Status = OK
}

func (gs *GameState) applyForageOrder(order *Order) {
	if gs.getUnit(order) == nil {
		return
	}

	hex := gs.Hexes[order.Coords]
	if hex.Terrain != FIELD || hex.Resources == 0 {
		order.Status = CANNOT_FORAGE
		return
	}

	hex.Resources--
	gs.PlayerResources[order.Player]++

	order.Status = OK
}

func (gs *GameState) applySpawnOrder(order *Order) {
	if gs.getUnit(order) == nil {
		return
	}
	if gs.targetIsBlocked(order) {
		return
	}
	if !gs.tryToPay(order, BEE_COST) {
		return
	}

	bee := &Entity{Type: BEE, HP: INIT_BEE_HP, Player: order.Player}
	gs.Hexes[order.Target()].Entity = bee

	order.Status = OK
}

type CoordsEntity struct {
	Coords Coords
	Entity *Entity
}

func (gs *GameState) Hives() []CoordsEntity {
	hives := []CoordsEntity{}

	for coords, hex := range gs.Hexes {
		if hex.Entity != nil && hex.Entity.Type == HIVE {
			hives = append(hives, CoordsEntity{coords, hex.Entity})
		}
	}

	return hives
}

func (gs *GameState) updateInfluence() {

	hives := gs.Hives()

	for coords, hex := range gs.Hexes {
		minDist := math.MaxUint32
		closestPlayers := make(map[int]bool)
		previousInfluence := hex.Influence

		for _, hive := range hives {
			dist := coords.Distance(hive.Coords)
			if dist > HIVE_FIELD_OF_VIEW {
				continue
			}

			if dist < minDist {
				minDist = dist
				clear(closestPlayers)
			}

			if dist <= minDist {
				closestPlayers[hive.Entity.Player] = true
			}
		}

		if len(closestPlayers) == 1 {
			for player := range closestPlayers {
				hex.Influence = player
			}
		} else {
			hex.Influence = -1
		}

		if (hex.Influence != previousInfluence) {
			gs.LastInfluenceChange = gs.Turn
		}
	}
}

func (gs *GameState) checkEndGame() {

	// No influence change in a while

	if gs.Turn - gs.LastInfluenceChange > INFLUENCE_TIMEOUT {
		gs.GameOver = true
		return
	}

	// Count influenced cells and hives

	influenceCounts := make([]int, gs.NumPlayers)
	hiveCounts := make([]int, gs.NumPlayers)

	for _, hex := range gs.Hexes {
		if hex.Influence >= 0 {
			influenceCounts[hex.Influence]++
		}
		if hex.Entity != nil && hex.Entity.Type == HIVE {
			hiveCounts[hex.Entity.Player]++
		}
	}

	// If a single player has hives, they win

	playersWithHives := 0
	for _, count := range hiveCounts {
		if count > 0 {
			playersWithHives++
		}
	}

	if playersWithHives == 1 {
		for player, count := range hiveCounts {
			if count > 0 {
				gs.Winners = append(gs.Winners, player)
				break
			}
		}
		gs.GameOver = true
		return
	}

	// Check if anyone has more than half the map influenced

	maxInfluence := slices.Max(influenceCounts)
	if maxInfluence <= len(gs.Hexes) / 2 {
		return
	}

	for player, influence := range influenceCounts {
		if influence == maxInfluence {
			gs.Winners = append(gs.Winners, player)
		}
	}

	gs.GameOver = true
}

func (gs *GameState) isVisibleBy(coords Coords, player int) bool {
	for hcoords, hex := range gs.Hexes {
		if hex.Entity != nil && hex.Entity.Player == player {
			if hcoords.Distance(coords) <= HIVE_FIELD_OF_VIEW {
				return true
			}
		}
	}
	return false
}

func (gs *GameState) PlayerView(player int) *GameState {
	view := &GameState{
		NumPlayers:          gs.NumPlayers,
		Turn:                gs.Turn,
		Hexes:               make(map[Coords]*Hex),
		LastInfluenceChange: gs.LastInfluenceChange,
		Winners:             gs.Winners,
		GameOver:            gs.GameOver,
	}

	for coords, hex := range gs.Hexes {
		if gs.isVisibleBy(coords, player) {
			view.Hexes[coords] = hex
		}
	}

	view.PlayerResources = []uint{gs.PlayerResources[player]}

	return view
}
