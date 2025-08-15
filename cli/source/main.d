import std.stdio;
import std.algorithm;
import std.range;
import std.json;

import game;
import terrain;
import order;
import serialization;

void main()
{
	auto map = loadMap("map.txt");
	auto game = new GameState(map[0], map[1], 3);

	auto move1 = new MoveOrder(game, player: 1, Coords(13, 3), Direction.W);
	auto move2 = new ForageOrder(game, player: 1, Coords(13, 1));

	game.processOrders([move1, move2]);
	writeln(move2.status);

	auto lol = serialize(game);
	auto foo = deserializeGameState(lol);
	writeln(foo);
}
