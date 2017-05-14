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
db = client["weather"]
all_collection = ["STM", "SS"]
gps_db = client["gps"]
gps_collection = ["account", "all_gps"]

class MainHandler(tornado.web.RequestHandler):
	""" MainHandler class. """

	def get(self):
		self.render("post_new_message_page.html")
	def post(self):
		collection_name = self.get_argument(name="collection_name")
		if collection_name not in all_collection:
			self.write("error collection name")
		elif collection_name == "STM":
			print self.request.body_arguments
			data = self.request.body_arguments
			try:
				db[collection_name].insert_one(data)
			except:
				self.write("error when insert message.")
			self.write("insert finished.")
		elif collection_name == "SS":
			data = self.request.body_arguments
			try:
				db[collection_name].insert_one(data)
			except:
				self.write("error when insert message.")
			self.write("insert finished.")
			

class SearchHandler(tornado.web.RequestHandler):
	"""PicHandler class. """
	def get(self):
		self.render("search_message_page.html")
	def post(self):
		data = ""
		collection_name = self.get_argument(name="collection_name")
		if collection_name not in all_collection:
			self.write("error collection name")
		elif collection_name == "STM":
			senser_id = self.get_argument(name="senser_id")	
			try:
				print senser_id
				data = db[collection_name].find_one({"senser_id":senser_id})
				pprint.pprint(data)
				self.write(str(data).replace(",","<br>"))
			except Exception,e:
				print e
				self.write("error when search message.")
		elif collection_name == "SS":
			senser_id = self.get_argument(name="senser_id")	
			try:
				print senser_id
				data = db[collection_name].find_one({"senser_id":senser_id})
				self.write(str(data).replace(",","<br>"))
			except Exception,e:
				print e
				self.write("error when insert message.")
			
class LoginHandler(tornado.web.RequestHandler):
	def post(self):
		c_name = 'account'
		data = self.request.body_arguments
		try:
			
			password = data["password"][0]
			username = data["username"][0]
			print username
			exist = gps_db[c_name].find_one({"username":username})
			if exist:
				if exist["password"] == password:
					self.write("success")
				else:
					self.write("fail")
			else:
				gps_db[c_name].insert_one({"username":username, "password":password})
				self.write("new account")
		except:
			self.write("error")

class RegisterHandler(tornado.web.RequestHandler):
	def post(self):
		c_name = 'account'
		data = self.request.body_arguments
		try:
			password = data["password"][0]
			username = data["username"][0]
			exist = gps_db[c_name].find_one({"username":username})
			if exist:
				self.write( "fail")
			else:
				gps_db[c_name].insert_one({"username":username, "password":password})
				self.write( "success")
		except:
			self.write( "error")
class PostGPSHandler(tornado.web.RequestHandler):
	def post(self):
		data = self.request.body_arguments
		try:
			password = data["password"][0]
			username = data["username"][0]
			new_data = {}
			new_data["username"] = username
			exist = gps_db["account"].find_one({"username":username})
			if exist:
				if exist["password"] == password:
					new_data["lat"] = float(data["lat"][0])
					new_data["lng"] = float(data["lng"][0])
					post_time = time.strftime("%Y%m%d%H%M%S", time.localtime()) 
					new_data["day"] = post_time[0:8]
					new_data["hour"] = post_time[8:14]
					gps_db["all_gps"].insert_one(new_data)
					self.write("success")
				else:
					self.write("wrong password")
			else:
				self.write( "noaccount")
		except Exception,e:
			print e
			self.write( "fail")
def my_order(x):
	return int(x["hour"])

class QueryGPSHandler(tornado.web.RequestHandler):
	def get(self):
		post_time = time.strftime("%Y%m%d%H%M%S", time.localtime()) 
		day = post_time[0:8]
		self.render("query_gps.html", day=day,username="dys")
	def post(self):
		data = self.request.body_arguments
		try:
			username = data["username"][0]
			day = data["day"][0]
			print username,day
		#	post_time = time.strftime("%Y%m%d%H%M%S", time.localtime()) 
		#	day = post_time[0:6]
			ans = []
			for a in gps_db["all_gps"].find({"username":username, "day":day}):
				a.pop("_id")
				ans.append(a)
			ans.sort(key=my_order)
			for a in ans:
				a["hour"] = "%s:%s:%s"%(a["hour"][0:2], a["hour"][2:4],a["hour"][4:6])
			self.render("query_map.html",jsonlist=ans,day=day,username=username)
		except Exception, e:
			print e
			self.write( "error" )



class DownloadHandler(tornado.web.RequestHandler):
	def get(self):
		buf_size=1024
		fileName = "test.apk"
		self.set_header('Content-Type','application/vnd.android.package-archive')
		self.set_header('Content-Disposition', 'attachment; filename=' + fileName)
		with open("download/"+fileName, 'r') as f:
			while True:
				data = f.read(buf_size)
				if not data:
					break
				self.write(data)
		self.finish()


class PicHandler(tornado.web.RequestHandler):
        def get(self):
		self.render("b2.html")



def make_app():
	settings = {
	"static_path": os.path.join(os.path.dirname(__file__), "static"),
	"autoreload":True,
	"debug":True,
	}
	return tornado.web.Application([
		(r"/post", MainHandler),
		(r"/search", SearchHandler),
		(r"/register", RegisterHandler),
		(r"/login", LoginHandler),
		(r"/querygps", QueryGPSHandler),
		(r"/postgps", PostGPSHandler),
		(r"/download", DownloadHandler),
		(r"/b2", PicHandler),
	],
	**settings)

if __name__ == "__main__":
	app = make_app()
	app.listen(2333)
	tornado.ioloop.IOLoop.current().start()

