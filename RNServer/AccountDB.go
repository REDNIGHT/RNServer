package RNServer

import (
	"errors"
	"strings"

	"../RNCore"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"
)

type AccountDB struct {
	RNCore.MongoDB

	inFind chan *FindAccount
}

type AccountData struct {
	Id        string
	LoginType string //手机 邮箱 设备id...
	Account   string
	Password  string
	Date      string //注册日期

	//AccountType string //普通用户 GM用户...
}

func NewAccountDB(name, url, user, pass, db, c string) *AccountDB {
	return &AccountDB{RNCore.NewMongoDB(name, url, user, pass, db, c), make(chan *FindAccount, RNCore.InChanLen)}
}

/*
func (this *ex)example() {
	ldb := NewAccountDB("", "", "", "")

	cb := func(ad *AccountData, err error) {
		if err != nil {
			return
		}
		this.SendMessage(func(_this RNCore.IMessage) {
			_ = ad.Account
		})
	}
	ldb.inFind <- &FindAccount{"acc", cb}
}
*/

func (this *AccountDB) Run() {
	defer this.CatchPanic()

	for {
		this.InTotal++

		//
		select {
		case i := <-this.inFind:
			is := make([]*FindAccount, len(this.inFind)+1)
			is[0] = i
			index := 1
			for i = range this.inFind {
				is[index] = i
				index++
			}
			this.find(is...)

			//
		case f := <-this.InCall():
			f(this)

		case f := <-this.InMessage():
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

type FindAccount struct {
	Account string
	CB      func(*AccountData, error)
}

func (this *AccountDB) Find(Account string, CB func(*AccountData, error)) {
	this.inFind <- &FindAccount{Account, CB}
}
func (this *AccountDB) find(is ...*FindAccount) {
	//todo...
	//测试是否返回一样个数的结果 会出现某项数据不存在而返回少量一个的情况

	arr := make([]bson.M, len(is))
	for index, i := range is {
		arr[index] = bson.M{"Account": i.Account}
	}
	query := bson.M{"$or": arr}

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

func (this *AccountDB) Insert(i *InsertAccount) {
	err := this.Collection.Insert(i.AD)
	i.CB(err)
}

type UpdatePassword struct {
	Account string
	//OldPassword string
	Password string
	CB       func(error)
}

func (this *AccountDB) UpdatePassword(i *UpdatePassword) {
	err := this.Collection.Update(bson.M{"Account": i.Account}, bson.M{"$set": bson.M{"Password": i.Password}})
	i.CB(err)
}

//
func (this *AccountDB) GetStateWarning(stateWarning func(name, warning string)) {
	this.TestChanOverload(stateWarning, "inFind", len(this.inFind))

	this.MongoDB.MNode.GetStateWarning(stateWarning)
}
