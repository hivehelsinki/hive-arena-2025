import std.random;
import std.conv;
import std.stdio;
import std.exception;
import std.datetime.systime;
import std.format;
import std.file;
import std.regex;
import std.algorithm;
import std.array;
import std.typecons;

import vibe.vibe;

import game;
import terrain;

const MAP_DIR = "maps";

alias GameID = uint;
alias Token = string;

struct GameResponse
{
	struct Cell
	{
		uint row, col;
		@byName Terrain terrain;

		@embedNullable Nullable!uint resources;
		@embedNullable Nullable!PlayerID influence;
		@embedNullable Nullable!Entity entity;
	}

	uint turn;

	Cell[] map;
	uint[] playerResources;
	uint lastInfluenceChange;

	PlayerID[] winners;
	bool gameOver;

	this(const GameState state, Nullable!PlayerID player = Nullable!PlayerID.init)
	{
		turn = state.turn;
		playerResources = state.playerResources.dup;
		lastInfluenceChange = state.lastInfluenceChange;

		foreach(coords, terrain; state.staticMap)
		{
			if (!player.isNull && !state.isVisibleBy(coords, player.get))
				continue;

			auto cell = Cell(
				coords.row,
				coords.col,
				terrain
			);

			if (terrain == Terrain.FIELD)
				cell.resources = state.mapResources[coords];

			if (coords in state.influence)
				cell.influence = state.influence[coords];

			if (coords in state.entities)
				cell.entity = cast(Entity) state.entities[coords];

			map ~= cell;
		}

		winners = state.winners.dup;
		gameOver = state.gameOver;
	}
}

class Player
{
	PlayerID id;
	string name;
	Token token;
}

class Game
{
	GameID id;
	string map;

	SysTime createdDate;

	Token adminToken;
	Token[] playerTokens;

	Player[] players;

	GameState state;

	static Token[] generateTokens(int count)
	{
		bool[Token] tokens;

		while (tokens.length < count)
		{
			auto token = format("%x", uniform!ulong);
			tokens[token] = true;
		}

		return tokens.keys;
	}

	this(GameID id, int numPlayers, MapData map)
	{
		this.id = id;
		this.map = map.name;

		createdDate = Clock.currTime;

		auto tokens = generateTokens(numPlayers + 1);
		adminToken = tokens[0];
		playerTokens = tokens[1 .. $];

		state = new GameState(map, numPlayers);
	}

	Player addPlayer(string name)
	{
		if (full)
			throw new Exception("Game is full");

		auto player = new Player();
		player.id = cast(PlayerID) players.length;
		player.name = name;
		player.token = playerTokens[player.id];

		players ~= player;
		return player;
	}

	bool full()
	{
		return players.length == state.numPlayers;
	}

	Json fullState()
	{
		return GameResponse(state).serializeToJson;
	}

	Json playerView(Token token)
	{
		auto player = cast(PlayerID) playerTokens.countUntil(token);
		return GameResponse(state, nullable(player)).serializeToJson;
	}
}

class Server
{
	MapData[string] maps;
	Game[GameID] games;

	this(ushort port)
	{
		loadMaps();

		auto router = new URLRouter;
		router.registerWebInterface(this);

		auto settings = new HTTPServerSettings();
		settings.port = port;

		listenHTTP(settings, router);
	}

	private void loadMaps()
	{
		foreach (path; dirEntries(MAP_DIR, SpanMode.shallow))
		{
			auto name = path.name.matchFirst(r"/(\w+)\.txt")[1];
			auto map = loadMap(path);

			map.name = name;
			maps[name] = map;
		}

		logInfo("Loaded maps: " ~ maps.keys.join(", "));
	}

	struct NewGameResponse
	{
		GameID id;
		PlayerID numPlayers;
		string map;
		SysTime createdDate;
		Token adminToken;
	}

	Json getNewgame(int players, string map)
	{
		GameID id;
		do { id = uniform!GameID; } while (id in games);

		if (map !in maps)
		{
			status(HTTPStatus.badRequest);
			return Json("Unknown map: " ~ map);
		}

		if (!GameState.validNumPlayers(players))
		{
			status(HTTPStatus.badRequest);
			return Json("Invalid player count: " ~ players.to!string);
		}

		Game game = new Game(id, players, maps[map]);
		games[id] = game;

		logInfo("Created game %d (%s, %d players)", id, map, players);
		setTimer(5.minutes, () => removeIfNotStarted(id));

		return NewGameResponse(
			id: game.id,
			numPlayers: game.state.numPlayers,
			map: game.map,
			createdDate: game.createdDate,
			adminToken: game.adminToken
		).serializeToJson;
	}

	private void removeIfNotStarted(GameID id)
	{
		if (!games[id].full)
		{
			games.remove(id);
			logInfo("Removed game %d because of timeout", id);
		}
	}

	struct StatusResponse
	{
		GameID id;
		PlayerID numPlayers;
		PlayerID playersJoined;
		string map;
		SysTime createdDate;
	}

	Json getStatus()
	{
		return games.values.map!(game => StatusResponse(
			id: game.id,
			numPlayers: game.state.numPlayers,
			playersJoined: cast(PlayerID) game.players.length,
			map: game.map,
			createdDate: game.createdDate
		)).array.serializeToJson;
	}

	struct JoinResponse
	{
		PlayerID id;
		Token token;
	}

	Json getJoin(GameID id, string name)
	{
		if (id !in games)
		{
			status(HTTPStatus.badRequest);
			return Json("Invalid game id: " ~ id.to!string);
		}

		if (name == "")
		{
			status(HTTPStatus.badRequest);
			return Json("Invalid name");
		}

		auto game = games[id];
		if (game.full)
		{
			status(HTTPStatus.badRequest);
			return Json("Game is full");
		}

		auto player = game.addPlayer(name);
		logInfo("Player %s joined game %d (#%d)", player.name, game.id, player.id);

		return JoinResponse(player.id, player.token).serializeToJson;
	}

	Json getGame(GameID id, Token token)
	{
		if (id !in games)
		{
			status(HTTPStatus.badRequest);
			return Json("Invalid game id: " ~ id.to!string);
		}

		auto game = games[id];
		if (!game.full)
		{
			status(HTTPStatus.badRequest);
			return Json("Game has not started");
		}

		if (token == game.adminToken)
			return game.fullState;

		if (game.playerTokens.canFind(token))
			return game.playerView(token);

		status(HTTPStatus.badRequest);
		return Json("Invalid token");
	}
}

void main()
{
	setLogFormat(FileLogger.Format.threadTime, FileLogger.Format.threadTime);

	auto server = new Server(8000);
	runApplication();

	writeln("Are we there yet?");
}
