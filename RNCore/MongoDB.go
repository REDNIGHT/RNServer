package RNCore

import (
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	MNode

	Session    *mgo.Session
	Collection *mgo.Collection
}

func NewMongoDB(name, url, user, pass, db, c string, indexKeys ...string) MongoDB {
	mdb := MongoDB{NewMNode(name), nil, nil}

	session, err := mgo.Dial(url)
	if err != nil {
		mdb.Panic(err, "err != nil")
	}
	session.SetMode(mgo.Monotonic, true)
	mdb.Session = session

	//
	DB := session.DB(db)
	err = DB.Login(user, pass)
	if err != nil {
		mdb.Panic(err, "err != nil")
	}

	//
	mdb.Collection = DB.C(c)
	err = mdb.Collection.EnsureIndexKey(indexKeys...)
	if err != nil {
		mdb.Panic(err, "err != nil")
	}

	return mdb
}
func (this *MongoDB) Close() {
	this.MNode.Close()

	this.Session.Close()
}
