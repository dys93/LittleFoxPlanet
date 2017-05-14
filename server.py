# coding=utf-8
"""test tornado."""
import json
import tornado.ioloop
import tornado.web
import os
import time
import pymongo
from pymongo import MongoClient
import pprint
client = MongoClient()
db = ["fox"]
all_collection = ["user", "account", "weather", "wager"]

class IndexHandler(tornado.web.RequestHandler):
	""" MainHandler class. """
	def get(self):
		self.render("index.html")
		
class LoginHandler(tornado.web.RequestHandler):
	def get(self):
		self.render("login.html")
	def post(self):
		c_name = 'user'
		data = self.request.body_arguments
		try:
			username = "initTestUser"
			password = "btwwtf"
			exist = db[c_name].find_one({"username":username})
			if exist:
				if exist["password"] == password:
					self.render("stake.html")
				else:
					self.write("Error Username or Password")
			else:
				gps_db[c_name].insert_one({"username":username, "password":password})
				self.write("new account")
				self.render("stake.html")		
		except:
			self.write("error")

class BetHandler(tornado.web.RequestHandler):
	def get(self):
		self.render("stake.html")	
	def post(self):
		c_name = 'account'
		data = self.request.body_arguments
		try:
			username = "initTestUser"
			password = "btwwtf"
			coin = int(data["coin"][0])
			cow = int(data["cow"][0])
			day = int(data["day"][0])
			weather = data["weather"]
			exist = db["user"].find_one({"username":username})
			if exist and exist["password"] == password:
				account = db["account"].find_one({"username":username})
				if account and account["coin"]>=coin and account["cow"]>=cow:
					account["coin"]-=coin
					account["cow"]-=cow
					db["account"].save(account)
					bet = {"username":username, 
					            "coin": coin,
					            "cow": cow,
					            "day": day,
					            "weather": weather,
					            "status": "init" }
				#????
				self.render("win.html")
			else:
				self.write("Error username")
		except:
			self.write("error")

class WinHandler(tornado.web.RequestHandler):
	""" MainHandler class. """
	def get(self):
		self.render("win.html")

class LoseHandler(tornado.web.RequestHandler):
	""" MainHandler class. """
	def get(self):
		self.render("lose.html")


class RankHandler(tornado.web.RequestHandler):
	""" MainHandler class. """
	def get(self):
		self.render("rank.html")

class EndHandler(tornado.web.RequestHandler):
	""" MainHandler class. """
	def get(self):
		self.render("end.html")

def make_app():
	settings = {
	"static_path": os.path.join(os.path.dirname(__file__), "static"),
	"autoreload":True,
	"debug":True,
	"template_path": os.path.join(os.path.dirname(__file__),"static")
	}
	return tornado.web.Application([
		(r"/index.html", IndexHandler),
		(r"/login.html", LoginHandler),
		(r"/stake.html", BetHandler),
		(r"/win.html", WinHandler),
		(r"/lose.html", LoseHandler),
		(r"/rank.html", RankHandler),
		(r"/end.html", EndHandler),
	],
	**settings)

if __name__ == "__main__":
	app = make_app()
	app.listen(2333)
	tornado.ioloop.IOLoop.current().start()

