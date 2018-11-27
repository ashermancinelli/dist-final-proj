



from socket import *
from threading import Thread
from concurrent.futures import ProcessPoolExecutor as Pool
import jobs

tpool = Pool(5)

def server(addr):
    sock = socket(AF_INET, SOCK_STREAM)
    sock.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
    sock.bind(addr)
    sock.listen(5)
    while True:
        client, ad = sock.accept()
        print('Connection to ', ad)
        Thread(target=handler, args=(client,), daemon=True).start()

def handler(client):
    while True:
        req = client.recv(100)
        if not req:
            break
        n = int(req)

        future = tpool.submit(jobs.job, n)
        result = future.result()

        client.send(f'response: {result}\n'.encode('ascii'))
    print('Closed\n')

server(('', 8080))
