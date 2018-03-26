package RNCore

import (
//"fmt"
)

//
type Node struct {
	MinNode
	MessageNode
}

func NewNode(name string) Node { return Node{NewMinNode(name), NewMessageNode()} }

func (this *Node) Name() string      { return this.MinNode.Name() }
func (this *Node) Type_Name() string { return this.MinNode.Type_Name() }
