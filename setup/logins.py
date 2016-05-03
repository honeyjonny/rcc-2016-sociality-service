from client import WebClient, serv_url
import tornado.ioloop
import tornado.gen

bots = []
POLLING = 500

def read_bots():	
	print("reading bots")
	with open("../bot_credentials.txt", "r") as cred:
		for line in cred:
			if "_session" in line or "True" in line:
				continue
			else:
				uname, passwd = line.strip().replace("(","").replace(")","").replace("'","").split(",")
				bots.append( (uname.strip(), passwd.strip() )  )

async def set_sessions():
	print("\n\n Reload sessions start \n\n")
	if len(bots) == 0:
		read_bots()

	if len(bots) > 0:
		print("login bots")
		for bot in bots:
			uname, passwd = bot[0], bot[1]
			print(uname, passwd)
			try:
				wc = WebClient(uname, passwd, serv_url)
				res = await wc.only_login()
				print(res)
			except Exception as e:
				print(e)
				continue


async def test():
	wc = WebClient("hj", "hj", serv_url)
	res = await wc.only_login()
	print(res)


async def reload_loop():
	print("Start reload loop")
	while True:
		await set_sessions()
		print("\n\n Reload sessions end \n\n,\nplan new reload")


if __name__ == '__main__':

	tornado.ioloop.IOLoop.current().spawn_callback(reload_loop)
	tornado.ioloop.IOLoop.current().start()