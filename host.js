


// https://gist.github.com/creationix/707146
net = require('net');
var clients = [];

net.createServer(sock => {
    sock.name = sock.remoteAddress + ":" + sock.remotePort;
    clients.push(sock);
    broadcast("join;" + sock.name + "\n", sock);

    sock.on('data', data => {
        broadcast(data, sock);
    });

    function broadcast(msg, sender) {
        clients.forEach(c => {
            if (c === sender) return;
            c.write(msg);
        })
        process.stdout.write(msg);
    }

}).listen(2007);

console.log("Host is running at localhost:2007\n");