package main


import (
	"./parent"
	"./child"
)

func main() {
	p := parent.NewParent()
	c := p.CreateNewChild()
	c.PrintParentMessage()
}