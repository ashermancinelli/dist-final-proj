

# Welcome to GoWar!!

Text-based battle royale game. Currently configured to run on local network, however the host could easily be port forewarded to be served publicly for online play. Each player seeks to attack other players and do the most total damage in the game. The last player alive may not be the winner. 

### Setup:
 `git clone https://github.com/ashermancinelli/dist-final-proj.git`
 
 `cd dist-final-proj`
 
 `. run`

### To play:

To play goWar, watch the console for any damage taken. 
To attack someone, run `attack [PLAYER NAME]` in the client console.

To list the players currently on the server, run `list` in the client console.

To see your score, run `score` in the client console.

## Technical details:

goWar Battle Royale is played by 2 or more players. Each player connects to a host system running `host.js`. All players run the `client.go` program where the game will prompt to enter a 3 character name. This is automated by running `. run`, which will begin the `client.go` program and configure everything for you. After name is entered, wait until all players are in lobby.

Communication on server takes the form:


`type;data1;data2;data3...`. 
For example, 

- `start/stop;`
- `attack;recieverName;attackerName`
- `meta;meta details`


For class demonstration, the `run` script will be configured for the local network. Otherwise, the client program may be ran with `go run client.go -host [HOST NAME] -port :[PORT NAME]` 

for host run `node host.js` to let port be set to the default of `2007`. Otherwise, run `node host.js port=[PORT NUMBER]`
