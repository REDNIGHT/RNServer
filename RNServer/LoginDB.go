package RNServer

import (
	"RNCore"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LoginDB struct {
	RNCore.Node

	url string
	db  string
	c   string
	//user, pass string

	InFind           chan *AccountFindData
	InInsert         chan *AccountInsertData
	InUpdatePassword chan *InUpdatePassword

	session    *mgo.Session
	collection *mgo.Collection
}

type AccountFindData struct {
	Account string
	CB      func(*AccountData, error)
}

type AccountInsertData struct {
	AD *AccountData
	CB func(error)
}

type InUpdatePassword struct {
	Account string
	//OldPassword string
	Password string
	CB       func(error)
}

type AccountData struct {
	ID       string
	Account  string
	Password string
	Type     string
}

func NewLoginDB(name string, url, db, c string) *LoginDB {
	return &LoginDB{RNCore.NewNode(name), url, db, c, make(chan *AccountFindData, RNCore.InChanLen), make(chan *AccountInsertData, RNCore.InChanLen), make(chan *InUpdatePassword, RNCore.InChanLen), nil, nil}
}

func (this *LoginDB) Run() {

	//
	session, err := mgo.Dial(this.url)
	if err != nil {
		this.Panic(err.Error())
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	this.session = session

	this.collection = session.DB(this.db).C(this.c)
	this.collection.EnsureIndexKey("Account", "ID")

	//
	for {
		this.InTotal++

		//
		select {
		case i := <-this.InFind:
			this.find(i)
		case i := <-this.InInsert:
			this.insert(i)
		case i := <-this.InUpdatePassword:
			this.updatePassword(i)

		case f := <-this.InMessage():
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

//
func (this *LoginDB) find(f *AccountFindData) {
	go func() {
		result := &AccountData{}
		err := this.collection.Find(bson.M{"Account": f.Account}).One(result)
		f.CB(result, err)
	}()
}

func (this *LoginDB) insert(f *AccountInsertData) {
	go func() {
		err := this.collection.Insert(f.AD)
		f.CB(err)
	}()
}

func (this *LoginDB) updatePassword(f *InUpdatePassword) {
	go func() {
		err := this.collection.Update(bson.M{"Account": f.Account}, bson.M{"$set": bson.M{"Password": f.Password}})
		f.CB(err)
	}()
}

//
/*func (this *LoginDB) MFind(account string) (*AccountData, error) {
	result := &AccountData{}
	err := this.collection.Find(bson.M{"Account": account}).One(result)
	return result, err
}

func (this *LoginDB) MInsert(ad *AccountData) error {
	return this.collection.Insert(ad)
}*/

//
func (this *LoginDB) DebugChanState(chanOverload chan *RNCore.ChanOverload) {
	this.TestChanOverload(chanOverload, "InGet", len(this.InFind))
	this.TestChanOverload(chanOverload, "InInsert", len(this.InInsert))

	this.Node.DebugChanState(chanOverload)
}
