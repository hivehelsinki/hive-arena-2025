import std.algorithm;

import map;

class GameState
{
	ubyte numPlayers;
	Hex[Coords] hexes;
	uint[] flowers;

	static GameState spawn(const(Hex[Coords]) baseMap, ubyte numPlayers)
	{
		ubyte[] playerMapping;
		assert(numPlayers <= 6);

		switch(numPlayers)
		{
			case 1: playerMapping = [0, 1, 0, 0, 0, 0, 0]; break;
			case 2: playerMapping = [0, 1, 0, 0, 2, 0, 0]; break;
			case 3: playerMapping = [0, 1, 0, 2, 0, 3, 0]; break;
			case 4: playerMapping = [0, 0, 1, 2, 0, 3, 4]; break;
			case 5: playerMapping = [0, 1, 2, 3, 4, 5, 0]; break;
			case 6: playerMapping = [0, 1, 2, 3, 4, 5, 6]; break;
			default: throw new Exception("Invalid player count: " ~ numPlayers);
		}

		Hex[Coords] hexes;
		foreach (coords, hex; baseMap)
		{
			Hex newHex = hex;
			newHex.player = playerMapping[newHex.player];
			if (newHex.kind.among(HexKind.HIVE, HexKind.BEE) && newHex.player == 0)
				newHex.kind = HexKind.EMPTY;

			hexes[coords] = newHex;
		}

		auto state = new GameState;
		state.numPlayers = numPlayers;
		state.hexes = hexes;
		state.flowers = new uint[numPlayers + 1];

		return state;
	}
}
