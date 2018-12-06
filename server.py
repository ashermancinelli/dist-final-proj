

from socket import *
from threading import Thread
from concurrent.futures import ProcessPoolExecutor as Pool
import jobs
import time

tpool = Pool(10)

class City(object): 
	def __init__(self, init_health, port=8080):
		self.port = port
		self.H = init_health
		self.createSock()

	def createSock(self):
		self.sock = socket(AF_INET, SOCK_STREAM)
    		sock.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)

	def closeSock(self):
		self.sock.close()


	def attack(self, host, port):
		self.sock.connect((host, port))
		print(f'attacked {host} {port}')
		self.closeSock()
		self.createSock()

	def recover(self, steps):
		self.sock.bind(('', self.port))
		self.sock.listen(5)
		print('Waiting...')
		c, addr = self.sock.accept()
		print(f'attacked by {c}')
		self.H -= 5
		print(f'health is now {self.H}')
		self.closeSock()
		self.createSock()



def server(addr):
    sock = socket(AF_INET, SOCK_STREAM)
    sock.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
    sock.bind(addr)
    sock.listen(5)
    while True:
        client, ad = sock.accept()
        print('Recieved attack from  ', ad)
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

if __name__ == '__main__':
	c = City(100)

	server(('', 8080))
