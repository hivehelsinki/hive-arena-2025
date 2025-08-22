import std.random;
import std.conv;
import std.stdio;

import vibe.vibe;

alias GameID = uint;

class Game
{
	GameID id;
	int numPlayers;
	string map;

	this(GameID id, int numPlayers, string map)
	{
		this.id = id;
		this.numPlayers = numPlayers;
		this.map = map;
	}
}

class Server
{
	Game[GameID] games;

	this(ushort port)
	{
		auto router = new URLRouter;
		router.registerWebInterface(this);

		auto settings = new HTTPServerSettings();
		settings.port = port;

		listenHTTP(settings, router);
	}

	Json getNewgame(int players, string map)
	{
		GameID id;
		do { id = uniform!GameID; } while (id in games);

		auto game = new Game(id, players, map);
		games[id] = game;

		return game.serializeToJson;
	}
}

void main()
{
	auto server = new Server(8000);
	runApplication();

	writeln("Are we there yet?");
}
