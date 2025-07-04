import game;
import map;

class Order
{
	enum Status
	{
		PENDING,

		OUT_OF_BOUNDS,
		TARGET_OUT_OF_BOUNDS,
		BAD_UNIT,
		BAD_PLAYER,

		BLOCKED,
		ATTACKED,
		DESTROYED,

		OK
	}

	Status status;
	ubyte player;
	Coords coords;

	void validate(GameState state)
	{
		if (coords !in state.hexes)
			status = Status.OUT_OF_BOUNDS;

		checkUnitType(state);

		if (state.hexes[coords].player != player)
			status = Status.BAD_PLAYER;
	}

	void checkUnitType(GameState state)
	{
		if (state.hexes[coords].kind != HexKind.BEE)
			status = Status.BAD_UNIT;
	}

	abstract void apply(GameState state);
}

class TargetOrder : Order
{
	Direction dir;

	override void validate(GameState state)
	{
		super.validate(state);

		if (coords.neighbour(dir) !in state.hexes)
			status = Status.TARGET_OUT_OF_BOUNDS;
	}
}

class MoveOrder : TargetOrder
{

}

class AttackOrder : TargetOrder
{

}

class ForageOrder : TargetOrder
{

}

class BuildWallOrder : TargetOrder
{

}

class HiveOrder : Order
{

}

class SpawnOrder : Order
{
	override void checkUnitType(GameState state)
	{
		if (state.hexes[coords].kind != HexKind.HIVE)
			status = Status.BAD_UNIT;
	}
}
