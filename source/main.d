import std.stdio;

import map;
import game;

void main()
{
	auto map = loadMap("map.txt");
	auto game = GameState.spawn(map, 4);

	write(game.hexes.mapToString);
}
