package parent

import (
	"fmt"
	"child"
)

type Parent struct {
	message string
}

func (parent *Parent) PrintMessage() {
	fmt.Println(parent.message)
}

func (parent *Parent) CreateNewChild() *child.Child {
	return child.NewChild(parent)
}

func NewParent() *Parent {
	return &Parent{message: "Hello World"}
}
