

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
var port = 2007; 
var stdin = process.openStdin();
var started = false;

syncNames = () => 'meta;all players;' + names.join(';') + '\n';

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
        process.stdout.write('INTERNAL: ' + data.toString().trim() + ' FROM: ' + sock.name + '\n');
        d = data.toString().trim();

        // if a new name is added to the server, append to the 
        if (d.split(';')[0] === 'name' && !started) names.push(d.split(';')[1]);
        if (d.split(';')[0] === 'death') {
            list.splice( list.indexOf( d.split(';')[1] ), 1 );
            if (names.length < 2) {
                broadcast("stop\n", 'admin')
            }
        }

    });
}).listen(port);

console.log('Listening on port: ' + port + ' with debug mode: ' + debug);











