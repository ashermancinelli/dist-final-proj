

# Welcome to GoWar!!

## Game mode 1:
goWar Battle Royale.

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


`type;data1;data2;data3...`
`start/stop;`
`attack;recieverName;attackerName`
`meta;meta details`


See Tom or Asher in class. `host.js` 
if running host on computer at the front of the EJ 214 class the local inet = `10.200.100.69`

Terminal commands:
for host run `node host.js` to let port be set to the default of `2007`. Otherwise, run `node host.js port=[PORT NUMBER]`

for client run `go run client.go -host [HOST NAME] -port :[PORT NAME]`
