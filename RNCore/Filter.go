package RNCore

import (
	"time"
)

type IFilterName interface {
	FilterName() string
}

type Filter struct {
	Node

	filters []_FilterInfo

	In chan IFilterName

	out chan<- interface{}
}

type _FilterInfo struct {
	filterName string
	pickUp     bool
	out        chan<- interface{}
}

func NewFilter(name string) *Filter {
	return &Filter{Node: NewNode(name), filters: make([]_FilterInfo, 0), In: make(chan IFilterName, InChanCount)}
}

func (this *Filter) SetOut(out chan<- interface{}, node_chan_name string) {
	this.out = out

	//
	this.SetOutNodeInfos("out", node_chan_name)
}

func (this *Filter) AddFilter(filterName string, pickUp bool, out chan<- interface{}) {
	this.filters = append(this.filters, _FilterInfo{filterName, pickUp, out})
}

func (this *Filter) Run() {

	for {
		//
		select {

		case iFilter := <-this.In:
			pickUp := false
			for _, f := range this.filters {
				if iFilter.FilterName() == f.filterName {
					f.out <- iFilter

					if f.pickUp == true {
						pickUp = f.pickUp
					}
				}
			}

			if pickUp == false {
				this.out <- iFilter
			}

		case <-time.After(time.Second * StateTime):
			this.State()
			continue

		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

//
type _FilterStateInfo struct {
	StateInfo
	In int
}

func (this *Filter) OnState() IStateInfo {
	return &_FilterStateInfo{StateInfo{this}, len(this.In)}
}
