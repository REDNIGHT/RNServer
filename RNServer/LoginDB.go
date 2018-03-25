package RNServer

import (
	"RNCore"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LoginDB struct {
	RNCore.MongoDB

	//InFind           chan *FindAccount
	//InInsert         chan *InsertAccount
	//InUpdatePassword chan *UpdatePassword
}

type AccountData struct {
	ID       string
	Account  string
	Password string
	Type     string
}

func NewLoginDB(name, url, db, c string) *LoginDB {
	//ldb := &LoginDB{RNCore.NewMongoDB(name, url, db, c), make(chan *FindAccount, RNCore.InChanLen), make(chan *InsertAccount, RNCore.InChanLen), make(chan *UpdatePassword, RNCore.InChanLen)}
	ldb := &LoginDB{RNCore.NewMongoDB(name, url, db, c, "Account", "ID")}
	return ldb
}

/*
func (this *_)example() {
	ldb : LoginDB
	ldb.Find("acc", func(ad *AccountData, err error) {
		if err != nil {
			return
		}
		this.SendMessage(func(_this RNCore.IMessage) {
			_ = ad.Account
		})
	})
}
*/
func (this *LoginDB) Find(Account string, cb func(*AccountData, error)) {
	go func() {
		result := &AccountData{}
		err := this.Collection.Find(bson.M{"Account": Account}).One(result)
		cb(result, err)
	}()
}

func (this *LoginDB) Insert(ad *AccountData, cb func(error)) {
	go func() {
		err := this.Collection.Insert(ad)
		cb(err)
	}()
}

func (this *LoginDB) UpdatePassword(Account, Password string, cb func(error)) {
	go func() {
		err := this.Collection.Update(bson.M{"Account": Account}, bson.M{"$set": bson.M{"Password": Password}})
		cb(err)
	}()
}

/*
func (this *LoginDB) Run() {
	for {
		select {
		case i := <-this.InFind:
			this.find(i)
		case i := <-this.InInsert:
			this.insert(i)
		case i := <-this.InUpdatePassword:
			this.updatePassword(i)
		}
	}
}

type FindAccount struct {
	Account string
	CB      func(*AccountData, error)
}

func (this *LoginDB) find(i *FindAccount) {
	go func() {
		result := &AccountData{}
		err := this.Collection.Find(bson.M{"Account": i.Account}).One(result)
		i.CB(result, err)
	}()
}

type InsertAccount struct {
	AD *AccountData
	CB func(error)
}

func (this *LoginDB) insert(i *InsertAccount) {
	go func() {
		err := this.Collection.Insert(i.AD)
		i.CB(err)
	}()
}

type UpdatePassword struct {
	Account string
	//OldPassword string
	Password string
	CB       func(error)
}

func (this *LoginDB) updatePassword(i *UpdatePassword) {
	go func() {
		err := this.Collection.Update(bson.M{"Account": i.Account}, bson.M{"$set": bson.M{"Password": i.Password}})
		i.CB(err)
	}()
}
*/

//
/*func (this *LoginDB) MFind(account string) (*AccountData, error) {
	result := &AccountData{}
	err := this.Collection.Find(bson.M{"Account": account}).One(result)
	return result, err
}

func (this *LoginDB) MInsert(ad *AccountData) error {
	return this.Collection.Insert(ad)
}*/
