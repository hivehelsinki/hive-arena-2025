import std.stdio;

import game;
import map;
import order;

void main()
{
	auto map = loadMap("map.txt");
	auto game = GameState.spawn(map, 4);

	write(game.hexes.mapToString);

	auto foo = new MoveOrder();

	game.applyOrders([foo]);
}
