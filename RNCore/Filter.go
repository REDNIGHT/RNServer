package RNCore

import (
//"time"
)

type Filter struct {
	In func() (string, interface{})

	outFilters []_FilterInfo

	Out func(interface{})
}

type _FilterInfo struct {
	filterName string
	pickUp     bool
	out        func(interface{})
}

func NewFilter() *Filter {
	return &Filter{nil, make([]_FilterInfo, 0), nil}
}

func (this *Filter) OutAddFilter(filterName string, pickUp bool, out func(interface{})) {
	this.outFilters = append(this.outFilters, _FilterInfo{filterName, pickUp, out})
}

func (this *Filter) Go() {
	go func() {
		for {
			name, v := this.In()
			pickUp := false
			for _, f := range this.outFilters {
				if name == f.filterName {
					f.out(v)

					if f.pickUp == true {
						pickUp = f.pickUp
					}
				}
			}

			if pickUp == false {
				this.Out(v)
			}

		}
	}()
}
