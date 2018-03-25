package RNServer

import (
	"RNCore"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
	//"gopkg.in/mgo.v2"
)

type LoginDB struct {
	RNCore.MongoDB

	InFind           chan *FindAccount
	InInsert         chan *InsertAccount
	InUpdatePassword chan *UpdatePassword
}

type AccountData struct {
	ID       string
	Account  string
	Password string
	Type     string
}

func NewLoginDB(name, url, db, c string) *LoginDB {
	return &LoginDB{RNCore.NewMongoDB(name, url, db, c), make(chan *FindAccount, RNCore.InChanLen), make(chan *InsertAccount, RNCore.InChanLen), make(chan *UpdatePassword, RNCore.InChanLen)}
}

/*
func example() {
	ldb := NewLoginDB("", "", "", "")

	cb := func(ad *AccountData, err error) {
		if err != nil {
			return
		}
		this.SendMessage(func(_this RNCore.IMessage) {
			_ = ad.Account
		})
	}
	ldb.InFind <- &FindAccount{"acc", cb}
}
*/

func (this *LoginDB) Run() {
	for {
		select {
		case i := <-this.InFind:
			is := make([]*FindAccount, len(this.InFind)+1)
			is[0] = i
			index := 1
			for i = range this.InFind {
				is[index] = i
				index++
			}
			this.find(is...)
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

func (this *LoginDB) find(is ...*FindAccount) {
	//todo...
	//测试是否返回一样个数的结果 会出现某项数据不存在而返回少量一个的情况

	arr := make([]bson.M, len(is))
	for index, i := range is {
		arr[index] = bson.M{"Account": i.Account}
	}
	query := bson.M{"$and": arr}

	iter := this.Collection.Find(query).Iter()
	defer iter.Close()

	if iter.Err() != nil {
		this.Error(iter.Err().Error())
		return
	}

	index := 0
	err_str := "<can not find> "
	result := &AccountData{}
	for iter.Next(result) {
		for _, i := range is {
			if i.Account == result.Account {
				i.CB(result, nil)
			} else {
				i.Account = err_str + i.Account
			}
		}
		index++
	}

	//
	if index != len(is) {
		this.Error("find  index != len(is) index=%v  len(is)=%v", index, len(is))
		for _, i := range is {
			if strings.Index(i.Account, err_str) == 0 {
				i.CB(nil, errors.New(err_str+i.Account))
			}
		}
	}
}

type InsertAccount struct {
	AD *AccountData
	CB func(error)
}

func (this *LoginDB) insert(i *InsertAccount) {
	err := this.Collection.Insert(i.AD)
	i.CB(err)
}

type UpdatePassword struct {
	Account string
	//OldPassword string
	Password string
	CB       func(error)
}

func (this *LoginDB) updatePassword(i *UpdatePassword) {
	err := this.Collection.Update(bson.M{"Account": i.Account}, bson.M{"$set": bson.M{"Password": i.Password}})
	i.CB(err)
}
