package RNCore

import (
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	Node

	Session    *mgo.Session
	Collection *mgo.Collection
}

func NewMongoDB(name, url, db, c string) MongoDB {
	mdb := MongoDB{NewNode(name), nil, nil}

	session, err := mgo.Dial(url)
	if err != nil {
		panic(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)
	mdb.Session = session

	mdb.Collection = session.DB(db).C(c)
	//this.collection.EnsureIndexKey("Account", "ID")

	return mdb
}

func (this *MongoDB) Run() {

	//
	for {
		this.InTotal++

		//
		select {

		case f := <-this.InMessage():
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

func (this *MongoDB) Close() {
	this.Close()
	this.Session.Close()
}
