local http = require "http.request"
local json = require "lunajson"

local args = {...}
local host = args[1]

local headers, stream = http.new_from_uri("http://" .. host .. "/newgame?players=4&map=balanced"):go()
local info = json.decode(stream:get_body_as_string())

for i = 1,4 do
	os.execute(string.format("go run hive-arena/example-agent-go %s %s %s %s",
		host,
		info.id,
		"Team" .. i,
		i == 4 and "" or " >> /dev/null &"
	))
end
