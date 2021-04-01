package models

import "gopkg.in/mgo.v2/bson"

func DashPassword() (passwords []bson.M, err error) {
	coll := Session.DB(DataName).C("password")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$site", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}, {"$limit": 20}})

	err = pipe.All(&passwords)
	return passwords, err
}

func DashUrls() (urls []bson.M, err error) {
	coll := Session.DB(DataName).C("urls")
	err = coll.Find(nil).Limit(20).All(&urls)

	return urls, err
}

func DashIps() (evilIps []bson.M, err error) {
	coll := Session.DB(DataName).C("evil_ips")
	err = coll.Find(nil).Limit(20).All(&evilIps)

	return evilIps, err
}

func DashTotal() (totalRecord int, totalPassword int, err error) {
	coll := Session.DB(DataName).C("proxy_honeypot")
	totalRecord, err = coll.Find(nil).Count()
	collPassword := Session.DB(DataName).C("password")
	totalPassword, err = collPassword.Find(nil).Count()
	return totalRecord, totalPassword, err
}
