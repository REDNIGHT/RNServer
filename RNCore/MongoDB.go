package RNCore

import (
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	MinNode

	Session    *mgo.Session
	Collection *mgo.Collection
}

func NewMongoDB(name, url, user, pass, db, c string, indexKeys ...string) MongoDB {
	mdb := MongoDB{NewMinNode(name), nil, nil}

	session, err := mgo.Dial(url)
	if err != nil {
		mdb.Panic(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)
	mdb.Session = session

	//
	DB := session.DB(db)
	err = DB.Login(user, pass)
	if err != nil {
		mdb.Panic(err.Error())
	}

	//
	mdb.Collection = DB.C(c)
	err = mdb.Collection.EnsureIndexKey(indexKeys...)
	if err != nil {
		mdb.Panic(err.Error())
	}

	return mdb
}

func (this *MongoDB) Close() {
	this.Session.Close()
}
