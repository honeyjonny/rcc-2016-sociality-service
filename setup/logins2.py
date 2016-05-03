from client import WebClient, serv_url
import tornado.ioloop
from time import sleep

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


def plan_future(future):
	print("\n\n Reload sessions end \n\n,\nplan new reload")
	tornado.ioloop.IOLoop.current().add_future(set_sessions(), plan_future)
	return


if __name__ == '__main__':

	#tornado.ioloop.IOLoop.current().add_future(set_sessions(), plan_future)
	tornado.ioloop.IOLoop.current().run_sync(set_sessions)