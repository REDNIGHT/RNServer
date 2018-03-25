package RNCore

import (
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	MinNode

	Session    *mgo.Session
	Collection *mgo.Collection
}

func NewMongoDB(name, url, db, c string, indexKeys ...string) MongoDB {
	mdb := MongoDB{NewMinNode(name), nil, nil}

	session, err := mgo.Dial(url)
	if err != nil {
		mdb.Panic(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)
	mdb.Session = session

	mdb.Collection = session.DB(db).C(c)
	err = mdb.Collection.EnsureIndexKey(indexKeys...)
	if err != nil {
		mdb.Panic(err.Error())
	}

	return mdb
}

func (this *MongoDB) Close() {
	this.Close()
	this.Session.Close()
}
