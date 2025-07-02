# Hive Arena 2025

## Gamestate and orders

- map format and loading
- static part: empty spaces, fields and stones, hive and bees spawns
- dynamic part
	- hives
	- bees
	- wax walls
	- fields current amount
	- player resources
	- influence?
- stateless function that processes gamestate and orders, produces new game state

## Headless cli

- take in JSON, spit out JSON, that's all

## Server

- match creation
	- submit: number of players, map (custom? from list?)
	- random name/url
	- delete after 1 minute if not started
- player registration to match
	- get secret token
- gamestate
	- rate limit 2 per second
- submit orders
- save games to disk
	- one JSON per each state plus orders
- route to browse/access history?
- route for interactive viewer
