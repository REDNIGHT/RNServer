package RNCore

type root struct {
	name  string
	nodes []INode
}

var _root = &root{ServerName, make([]INode, 0)}

func Root() *root {
	return _root
}

func (this *root) Name() string { return this.name }

func (this *root) Add(node INode) {
	if this.Get(node.Name()) != nil {
		panic("Get(node.Name()) != nil  node.Name()=" + node.Name())
	}

	this.nodes = append(this.nodes, node)
}

func (this *root) Get(name string) INode {
	for i := 0; i < len(this.nodes); i++ {
		if this.nodes[i].Name() == name {
			return this.nodes[i]
		}
	}

	return nil
}

func (this *root) Init() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Init()

		n.Log("Init()")
	}
}

func (this *root) Register() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Register()

		n.Log("Register()")
	}
}

func (this *root) Run() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		go n.Run()

		n.Log("Run()")
	}
}

func (this *root) Close() {
	for i := len(this.nodes) - 1; i >= 0; i-- {

		n := this.nodes[i]
		n.Close()

		n.Log("Close()")
	}
}

func (this *root) Destroy() {
	for i := len(this.nodes) - 1; i >= 0; i-- {

		n := this.nodes[i]
		n.Destroy()

		n.Log("Destroy()")
	}
}

func destroy(n INode) {
	/*defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()*/

	n.Destroy()
}
