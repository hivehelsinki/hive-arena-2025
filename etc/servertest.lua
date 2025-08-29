local http = require "http.request"
local json = require "lunajson"

local host = "http://localhost:8000/"

local function req(route, ...)
	local headers, stream = http.new_from_uri(string.format(host .. route, ...)):go()
	return json.decode(stream:get_body_as_string())
end

local function start_game(players, map)

	local game = req("newgame?players=%d&map=%s", players, map)
	game.players = {}

	for i = 1, game.numPlayers do
		local player = req("join?id=%d&name=%s", game.id, "coolplayer" .. math.random(1000000))
		table.insert(game.players, player)
	end

	return game
end

local g = start_game(math.random(2,6), "balanced")

print(json.encode(g))
