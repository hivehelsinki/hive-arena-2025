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

class Order
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

	enum Status
	{
		PENDING,

		OUT_OF_BOUNDS,
		TARGET_OUT_OF_BOUNDS,
		BAD_UNIT,
		ENNEMY_UNIT,

		BLOCKED,
		ATTACKED,
		DESTROYED,

		OK
	}

	Status status;
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

	void validateOrder(Order order) const
	{
		if (order.coords !in hexes)
		{
			order.status = Order.Status.OUT_OF_BOUNDS;
			return;
		}

		Hex from = hexes[order.coords];

		if ((order.isBeeOrder && from.kind != HexKind.BEE) ||
			(order.isHiveOrder && from.kind != HexKind.HIVE))
		{
			order.status = Order.Status.BAD_UNIT;
			return;
		}

		if (from.player != order.player)
		{
			order.status = Order.Status.ENNEMY_UNIT;
			return;
		}

		if (order.hasTarget && order.coords.neighbour(order.dir) !in hexes)
		{
			order.status = Order.Status.TARGET_OUT_OF_BOUNDS;
			return;
		}
	}

	void applyOrders(Order[] orders)
	{
		orders.each!(a => validateOrder(a));
		auto valid = orders.filter!(a => a.status == Order.Status.PENDING);

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
			attack.status = Order.Status.OK;
		}

		// Invalidate orders to destroyed units

		foreach(order; valid)
			if (hexes[order.coords].hp == 0)
				order.status = Order.Status.DESTROYED;

		auto alive = valid.filter!(order => order.status != Order.Status.DESTROYED);

		// Movement

		solveMovements(alive.filter!(a => a.kind == Order.Kind.MOVE).array);

		// Spawns

		// All the other orders fail if the unit was attacked

		// Forage

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
