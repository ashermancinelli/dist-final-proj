

# Welcome to GoWar!!

## Game mode 1:
goWar Battle Royale.

### Setup:
goWar Battle Royale is played by 2 or more players. Each player connects to a hosting system running `host.js`. All players run the `client.go` program where the game will prompt to enter a 3 character name. After name is entered, wait until all players are in lobby. One VIP will enter the command goWarStart1. After this all players will recieve a start game notification and play begins.

### To play:
To play goWar, watch the console for any damage taken while typing in the name of an enemy to attack. 
---------------------------

## Technical details:

Communication on server takes the form:


`type;data1;data2;data3...`
`start/stop;`
`attack;recieverName;attackerName`
`meta;meta details`


See Tom or Asher in class. `host.js` 
if running host on computer at the front of the EJ 214 class the local inet = `10.200.100.69`

Terminal commands:
for host run `node host.js`

for client run `go run client.go -host 10.200.100.69 -port :2007`

if host is not being run on front of class machine change host ip names to whatever inet value host computer has.
