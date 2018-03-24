package RNServer

import (
	"RNCore"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LoginDB struct {
	RNCore.MongoDB

	InFind           chan *AccountFind
	InInsert         chan *AccountInsert
	InUpdatePassword chan *UpdatePassword
}

type AccountFind struct {
	Account string
	CB      func(*AccountData, error)
}

type AccountInsert struct {
	AD *AccountData
	CB func(error)
}

type UpdatePassword struct {
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

func NewLoginDB(name, url, db, c string) *LoginDB {
	return &LoginDB{RNCore.NewMongoDB(name, url, db, c), make(chan *AccountFind, RNCore.InChanLen), make(chan *AccountInsert, RNCore.InChanLen), make(chan *UpdatePassword, RNCore.InChanLen)}
}

func (this *LoginDB) Run() {

	this.Collection.EnsureIndexKey("Account", "ID")

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
func (this *LoginDB) find(i *AccountFind) {
	go func() {
		result := &AccountData{}
		err := this.Collection.Find(bson.M{"Account": i.Account}).One(result)
		i.CB(result, err)
	}()
}

func (this *LoginDB) insert(i *AccountInsert) {
	go func() {
		err := this.Collection.Insert(i.AD)
		i.CB(err)
	}()
}

func (this *LoginDB) updatePassword(i *UpdatePassword) {
	go func() {
		err := this.Collection.Update(bson.M{"Account": i.Account}, bson.M{"$set": bson.M{"Password": i.Password}})
		i.CB(err)
	}()
}

//
/*func (this *LoginDB) MFind(account string) (*AccountData, error) {
	result := &AccountData{}
	err := this.Collection.Find(bson.M{"Account": account}).One(result)
	return result, err
}

func (this *LoginDB) MInsert(ad *AccountData) error {
	return this.Collection.Insert(ad)
}*/

//
func (this *LoginDB) DebugChanState(chanOverload chan *RNCore.ChanOverload) {
	this.TestChanOverload(chanOverload, "InFind", len(this.InFind))
	this.TestChanOverload(chanOverload, "InInsert", len(this.InInsert))
	this.TestChanOverload(chanOverload, "InUpdatePassword", len(this.InUpdatePassword))

	this.Node.DebugChanState(chanOverload)
}
