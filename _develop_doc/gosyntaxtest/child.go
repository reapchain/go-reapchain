//When compile , error message : import cycle not allowed

// How to solve it.

// see this files.

// package dependency check command:  go list -f '{{join .Deps "\n"}}'

package child


import "Parent"

type Child struct {
	parent *Parent
}

func (child *Child) PrintParentMessage() {
	child.parent.PrintMessage()
}

func NewChild(parent *Parent) *Child {
	return &Child{parent: parent }
}