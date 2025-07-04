import std.stdio;
import std.algorithm;
import std.range;

import game;
import map;
import order;

void main()
{
	auto map = loadMap("map.txt");
	auto game = GameState.spawn(map, 4);

	write(game.hexes.mapToString);

//	write(game.hexes.byPair.filter!(a => a.value.kind == HexKind.BEE));

	auto foo = new MoveOrder();
	foo.player = 1;
	foo.coords = Coords(2,14);
	foo.dir = Direction.SE;

	auto bar = new AttackOrder();
	bar.player = 1;
	bar.coords = Coords(2,12);
	bar.dir = Direction.NE;


	game.applyOrders([foo, bar]);
	writeln(foo.status);

	write(game.hexes.mapToString);
	writeln(game.hexes[Coords(1,13)]);
}
