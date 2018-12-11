

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


net = require('net');
var clients = [];
var names = [];
var points;
var port = 2007; 
var stdin = process.openStdin();
var started = false;
var deadCount = 0;

syncNames = () => 'meta;all players;' + names.join(';') + '\n';

printScore = () => {
    var score = '\n\nScore:\n--------------------\n';
    for (var i = 0; i < names.length; i++) {
        score += names[i] + '\t\t' + points[i] + '\n\n';
    }
    return score + '\n\n\n';
}

process.argv.forEach( arg => {
    if (arg.includes('port')) {
        port = parseInt(arg.split('=')[1], 10);
    }
});

broadcast = (msg, sender) => {
    clients.forEach(c => {
        if (c === sender) return;
        c.write(msg);
    });
}

stdin.addListener('data', d => {
    d = d.toString().trim();
    if (d === 'names') {
        d = syncNames();
    } else if (d === 'start') {
        started = true;
        points = new Array(names.length).fill(0);
        broadcast(syncNames(), 'admin');
        d = 'start\n'
    }
    broadcast(d, 'admin');
    process.stdout.write('Sent admin command: ' + d + '\n');
});

net.createServer(sock => {
    sock.name = sock.remoteAddress + ":" + sock.remotePort;
    clients.push(sock);
    broadcast("join;" + sock.name + "\n", sock);

    sock.on('data', data => {
        broadcast(data, sock);
        process.stdout.write('INTERNAL:\t\t' + data.toString().trim() + '\t\tFROM: ' + sock.name + '\n');
        d = data.toString().trim();

        // if a new name is added to the server, append to the 
        if (d.split(';')[0] === 'name' && !started) {
            names.push(d.split(';')[1]);
            process.stdout.write('Received name from command: ' + d + '\n');
        } else if (d.split(';')[0] === 'death') {
            deadCount++;
            if (names.length - deadCount < 2) {
                broadcast("stop;" + printScore(), 'admin')
            }
        } else if (d.split(';')[0] === 'attack') {
            points[names.indexOf(d.split(';')[2])] += parseInt(d.split(';')[3]);
        }

    });
}).listen(port);

console.log('Listening on port: ' + port);











