

/*
Asher Mancinelli
Host program which takes some number of clients via TCP socket
and forms a 'room' where each message from each socket is 
transmitted down every other socket in the room.

In the game, each client will send message data about its
agent, and the host will do some light filtering and reflect
that message across all other sockets, and those individual
programs will interpret the game data.

Reference:
https://gist.github.com/creationix/707146

*/

// sets variables to be used
net = require('net');
var clients = [];
var names = [];
var points;
var port = 2007; 
var stdin = process.openStdin();
var started = false;
var deadCount = 0;
// ----------------------

// inline function to send the signal for all players
syncNames = () => 'meta;all players;' + names.join(';') + '\n';

// creates a scoreboard string which will be passed along to all
// clients when the game ends so that each client does not have to 
// keep track of the gloabal scoreboard
printScore = () => {
    var score = '\n\nScore:\n--------------------\n';
    for (var i = 0; i < names.length; i++) {
        score += names[i] + '\t\t' + points[i] + '\n\n';
    }
    var maxIdx = points.indexOf(Math.max.apply(null, points));
    var winner = names[maxIdx];
    score += 'Winner is: ' + winner + '!!!\n'
    return score + '\n\n\n';
}

// parses command line args for the 
// port to listen on
process.argv.forEach( arg => {
    if (arg.includes('port')) {
        port = parseInt(arg.split('=')[1], 10);
    }
});

// function to send a message from one socket 
// or generic sender ('admin' when sending from host)
// to all other sockets (clients)
broadcast = (msg, sender) => {
    clients.forEach(c => {
        if (c === sender) return;
        c.write(msg);
    });
}

// callback to be called whenver data can be pulled from stdin
stdin.addListener('data', d => {
    d = d.toString().trim();

    // sync the names of all the clients when the host 
    // calls `names`
    if (d === 'names') {
        d = syncNames();

    // when the host runs the `start` command, the host 
    // sends the signals to sync the names across all the clients,
    // and then sends the `start` command.
    } else if (d === 'start') {
        started = true;
        points = new Array(names.length).fill(0);
        broadcast(syncNames(), 'admin');
        d = 'start\n'
    }

    // otherwise, broadcast the admin's command as raw communication
    // to all the clients. 
    broadcast(d, 'admin');
    process.stdout.write('Sent admin command: ' + d + '\n');
});

// creates port to listen on, and accepts all the incoming
// requests for sockets and adds them to the communication room
net.createServer(sock => {

    sock.name = sock.remoteAddress + ":" + sock.remotePort;

    // adds the newly joined client to the list of all clients
    clients.push(sock);

    // signals a new player has joined to host and all clients. 
    broadcast("join;" + sock.name + "\n", sock);

    // sets event listener (callback) to run when data is recieved down 
    // any of the pipes. Processes input from the sockets and sometimes
    // broadcasts data back down the pipes in response to the given data
    sock.on('data', data => {

        // send data down pipe to all other clients
        broadcast(data, sock);

        // prints internally on the host when any command is recieved
        process.stdout.write('INTERNAL:\t\t' + data.toString().trim() + '\t\tFROM: ' + sock.name + '\n');
        d = data.toString().trim();

        // if a new name is added to the server, append to the list of names
        if (d.split(';')[0] === 'name' && !started) {
            names.push(d.split(';')[1]);
            process.stdout.write('Received name from command: ' + d + '\n');
        } else if (d.split(';')[0] === 'death') {
            // if someone dies, increment the dead counter
            deadCount++;

            // if there is one player left, end the game and print out the score
            if (names.length - deadCount < 2) {
                broadcast("stop;" + printScore(), 'admin')
            }
        } else if (d.split(';')[0] === 'attack') {

            // keep track of score when an attack is sent
            points[names.indexOf(d.split(';')[2])] += parseInt(d.split(';')[3]);
        }

    });

// listen on the given port (default 2007)
}).listen(port);

console.log('Listening on port: ' + port);
