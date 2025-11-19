local arena = require "arena"
local utils = require "utils"

local args = {...}

local host = args[1]
local gameid = args[2]
local name = args[3]

if not host or not gameid or not name then
	print "Usage: lua main.lua <host> <gameid> <name>"
	return
end

local function makeOrders(state, player)

	utils.printState(state)

	local orders = {}

	for coords,hex in pairs(state.hexes) do
		if hex.entity and hex.entity.type == "BEE" and hex.entity.player == player then

			local order = {
				type = "MOVE",
				coords = coords,
				direction = ({"NE", "E", "SE", "SW", "W", "NW"})[math.random(1,6)]
			}

			table.insert(orders, order)
		end
	end

	return orders
end

arena.runAgent(host, gameid, name, makeOrders)
