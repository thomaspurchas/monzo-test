package grapher

import "fmt"

type Node struct {
	id   int
	name string
}

func (n Node) ID() int {
	return n.id
}

func (n Node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.name)
}
