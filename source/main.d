import std.stdio;

import map;
import game;

void main()
{
	auto map = loadMap("map.txt");
	auto game = GameState.spawn(map, 3);

	write(game.hexes.mapToString);
}
