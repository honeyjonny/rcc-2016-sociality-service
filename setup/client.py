
import tornado.ioloop

from random import randrange

import random
import string

from loremipsum import *
from tornado.escape import *
from tornado.httpclient import *

FLAG = "VolgaCTF{__why_so_strange_request_statuses__}"

class WebClient(object):

	def __init__(self, username, password, url):
		self.username = username
		self.password = password
		self.url = url
		self.client = AsyncHTTPClient()
		self.headers = {
		"Content-Type":"application/x-www-form-urlencoded",
		"User-Agent":"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36"
		}

	def add_cookie(self, cookie):
		nheaders = dict(self.headers)
		nheaders["Cookie"] = cookie
		return nheaders

	def get_post_body(self):
		return "username=" + self.username + "&password=" + self.password

	async def register_usr(self):

		body = self.get_post_body()

		req = HTTPRequest(self.url+"/register", 
							method="POST", 
							headers=self.headers, 
							body=body, 
							follow_redirects=False)

		try:
			resp = await self.client.fetch(req)
		except HTTPError as err:
			resp = err.response

		return resp.code

		

	async def login_user(self):

		body = self.get_post_body()

		req = HTTPRequest(self.url+"/login", 							
							method="POST", 
							headers=self.headers, 
							body=body, 
							follow_redirects=False)

		try:
			resp = await self.client.fetch(req)
		except HTTPError as err:
			resp = err.response

		#print(resp.code, resp.body, resp.headers["Set-Cookie"])
		if resp.code == 302:
			return resp.headers["Set-Cookie"].split(";")[0]
		else:
			return None

	async def do_lorem_posts(self, cookie):

		sent_num = randrange(5,15)

		sents = get_sentences(sent_num)

		if self.username == "uuser173":
			sents.insert(4, FLAG)

		for s in sents:

			body = "content=" + s

			headers = self.add_cookie(cookie)

			req = HTTPRequest(self.url+"/posts", 							
					method="POST", 
					headers=headers, 
					body=body, 
					follow_redirects=False)

			try:
				resp = await self.client.fetch(req)
			except HTTPError as err:
				resp = err.response

			if resp.code == 303:
				continue
			else:
				raise Exception("some err in do lorem posts")

		return True


	async def start(self):
		stat = await self.register_usr()
		if stat == 303:
			cookie = await self.login_user()

			if cookie != None:
				print(cookie)
				res = await self.do_lorem_posts(cookie)
				print(res)

	async def only_login(self):
		cookie = await self.login_user()
		if cookie == None:
			return None
		else:
			return await self.login_user()


async def randomword(length):
   return ''.join(random.choice(string.ascii_lowercase) for i in range(length))


serv_url = "http://ur.2016.volgactf.ru:8080"

async def main():	
	ustart = 1
	uend = 317
	ucounter = ustart
	
	while ucounter < uend:

		username = "uuser" + str(ucounter)
		password = await randomword(7)

		ucounter += 1

		print( (username, password) )

		wc = WebClient(username, password, serv_url)
		wc.start()


if __name__ == '__main__':
	tornado.ioloop.IOLoop.current().run_sync(main)
