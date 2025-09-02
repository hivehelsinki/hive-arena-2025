import std.stdio;
import std.json;
import std.getopt;
import std.algorithm;
import std.array;

import game;
import terrain;
import order;
import serialization;

void process(JSONValue data)
{
	auto game = deserializeGameState(data["gamestate"]);

	if (data["orders"].array.length != game.numPlayers)
		throw new Exception("Player count mismatch");

	Order[][] orders;
	foreach (Player p; 0 .. game.numPlayers)
		orders ~= data["orders"][p].array.map!(o => deserializeOrder(o, p, game)).array;

	auto processed = game.processOrders(orders);

	JSONValue result;
	result["gamestate"] = serialize(game);
	result["processed"] = processed.map!(a => serialize(a)).array;

	writeln(result);
}

void main(string[] args)
{
	string mapPath;
	ubyte players;
	string toProcess;
	bool help;

	getopt(args,
		"map", &mapPath,
		"players", &players,
		"process", &toProcess,
		"help", &help
	);

	if (help)
	{
		writeln("Usage:");
		writeln(args[0], " --map=<map file> --players=<number of players to spawn>");
		writeln(args[0], " --process=<json formatted game state and orders>");
		writeln("Without arguments, process standard input");

		return;
	}

	if (mapPath.length != 0)
	{
		auto map = loadMap(mapPath);
		auto game = new GameState(map, players);

		writeln(serialize(game));
		return;
	}

	JSONValue data;

	if (toProcess.length != 0)
	{
		import std.file;

		auto txt = readText(toProcess);
		data = parseJSON(txt);
	}
	else
	{
		string[] input;
		foreach (string line; lines(stdin))
			input ~= line;

		string txt = input.join;
		data = parseJSON(txt);
	}

	try
	{
		process(data);
	}
	catch (Exception e)
	{
		stderr.writeln(e);
		stderr.writeln(data);
	}

}
