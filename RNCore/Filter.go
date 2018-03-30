package RNCore

//"time"

type Filter struct {
	In func() (string, interface{})

	outFilters []func(inName string, v interface{}) bool

	Out func(interface{})
}

func NewFilter() *Filter {
	return &Filter{nil, make([]func(string, interface{}) bool, 0), nil}
}

func (this *Filter) OutAddFilter(outFilters ...func(inName string, v interface{}) bool) {
	this.outFilters = append(this.outFilters, outFilters...)
}

func (this *Filter) Run() {
	defer CatchPanic()

	for {
		name, v := this.In()
		pickUp := false
		for _, outFilter := range this.outFilters {
			if outFilter(name, v) {
				pickUp = true
			}
		}

		if pickUp == false {
			this.Out(v)
		}
	}
}
