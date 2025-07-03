import std.algorithm;
import std.array;
import std.conv;

import map;
import utils;

const ubyte[HexKind] maxHP = [
	HexKind.HIVE: 12,
	HexKind.BEE: 2,
	HexKind.FIELD: 120,
	HexKind.WALL: 6,
	HexKind.EMPTY: 0,
	HexKind.ROCK: 0
];

struct Order
{
	enum Kind
	{
		MOVE,
		HIVE,
		BUILD,
		ATTACK,
		FORAGE,
		SPAWN
	}

	ubyte player;
	Kind kind;
	Coords coords;
	Direction dir;

	bool isBeeOrder()
	{
		return kind.among(
			Kind.MOVE,
			Kind.HIVE,
			Kind.BUILD,
			Kind.ATTACK,
			Kind.FORAGE
		) != 0;
	}

	bool isHiveOrder()
	{
		return kind == Kind.SPAWN;
	}

	bool hasTarget()
	{
		return kind.among(
			Kind.MOVE,
			Kind.BUILD,
			Kind.ATTACK
		) != 0;
	}

	bool isLocal()
	{
		return !hasTarget;
	}
}

enum OrderStatus
{
	OK,
	OUT_OF_BOUNDS,
	TARGET_OUT_OF_BOUNDS,
	BAD_UNIT,
	ENNEMY_UNIT
}

class GameState
{
	ubyte numPlayers;
	Hex[Coords] hexes;
	uint[] flowers;

	static GameState spawn(const Hex[Coords] baseMap, ubyte numPlayers)
	{
		ubyte[] playerMapping;
		assert(numPlayers <= 6);

		switch (numPlayers)
		{
			case 1: playerMapping = [0, 1, 0, 0, 0, 0, 0]; break;
			case 2: playerMapping = [0, 1, 0, 0, 2, 0, 0]; break;
			case 3: playerMapping = [0, 1, 0, 2, 0, 3, 0]; break;
			case 4: playerMapping = [0, 0, 1, 2, 0, 3, 4]; break;
			case 5: playerMapping = [0, 1, 2, 3, 4, 5, 0]; break;
			case 6: playerMapping = [0, 1, 2, 3, 4, 5, 6]; break;
			default: throw new Exception("Invalid player count: " ~ numPlayers.to!string);
		}

		Hex[Coords] hexes;
		foreach (coords, baseHex; baseMap)
		{
			Hex hex = baseHex;

			hex.player = playerMapping[hex.player];
			if (hex.kind.among(HexKind.HIVE, HexKind.BEE) && hex.player == 0)
				hex.kind = HexKind.EMPTY;
			hex.hp = maxHP[hex.kind];

			hexes[coords] = hex;
		}

		auto state = new GameState;
		state.numPlayers = numPlayers;
		state.hexes = hexes;
		state.flowers = new uint[numPlayers + 1];

		return state;
	}

	OrderStatus isValidOrder(Order order) const
	{
		if (order.coords !in hexes)
			return OrderStatus.OUT_OF_BOUNDS;

		Hex from = hexes[order.coords];

		if (order.isBeeOrder && from.kind != HexKind.BEE)
			return OrderStatus.BAD_UNIT;

		else if (order.isHiveOrder && from.kind != HexKind.HIVE)
			return OrderStatus.BAD_UNIT;

		if (from.player != order.player)
			return OrderStatus.ENNEMY_UNIT;

		if (order.hasTarget && order.coords.neighbour(order.dir) !in hexes)
			return OrderStatus.TARGET_OUT_OF_BOUNDS;

		return OrderStatus.OK;
	}

	void applyOrders(Order[] orders)
	{
		auto valid = orders.filter!(a => isValidOrder(a) == OrderStatus.OK);

		// Attacks first

		bool[Coords] wasAttacked;
		foreach (attack; valid.filter!(a => a.kind == Order.Kind.ATTACK))
		{
			auto target = attack.coords.neighbour(attack.dir);
			hexes[target].hp--;

			if (hexes[target].hp == 0)
			{
				hexes[target].kind = HexKind.EMPTY;
				hexes[target].player = 0;
			}

			wasAttacked[target] = true;
		}

		// Invalidate orders to destroyed units

		auto alive = valid.filter!(order => hexes[order.coords].hp != 0);

		// Movement

		solveMovements(alive.filter!(a => a.kind == Order.Kind.MOVE).array);

		// Forage
		// Spawns
		// Hive starts
		// Builds
	}

	void solveMovements(Order[] moves)
	{
		// Moving into structures fail
		// Direct swaps fail
		// Contested destinations: one randomly success, the others fail
		// Cycles succeed
		// Chains succeed
	}
}
