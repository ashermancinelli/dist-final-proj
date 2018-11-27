

from socket import *

def server(addr):
    sock = socket(AF_INET, SOCK_STREAM)
    sock.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
    sock.bind(addr)
    sock.listen(5)
    while True:
        client, ad = sock.accept()
        print('Connection to ', ad)
        handler(client)

def handler(client):
    while True:
        req = client.recv(100)
        if not req:
            break
        client.send(f'responding to {req}\n'.encode('ascii'))
    print('Closed\n')

server(('', 8080))
