package grapher

import (
	"github.com/gonum/graph/encoding/dot"
	"github.com/gonum/graph/simple"
)

type Graph struct {
	*simple.DirectedGraph
	name string
}

func newGraph() *Graph {
	dg := simple.NewDirectedGraph(1.0, 1.0)
	g := &Graph{dg, ""}
	return g
}

func (g *Graph) DOTID() string {
	return g.name
}

type gAttributer struct{}
type nAttributer struct{}
type eAttributer struct{}

func (gAttributer) DOTAttributes() []dot.Attribute {
	dir := dot.Attribute{"rankdir", "\"LR\""}
	sep := dot.Attribute{"ranksep", "4.0"}
	a := []dot.Attribute{dir, sep}
	return a
}

func (nAttributer) DOTAttributes() []dot.Attribute {
	shape := dot.Attribute{"shape", "box"}
	a := []dot.Attribute{shape}
	return a
}

func (eAttributer) DOTAttributes() []dot.Attribute {
	head := dot.Attribute{"headport", "w"}
	tail := dot.Attribute{"tailport", "e"}
	return []dot.Attribute{head, tail}
}

func (g *Graph) DOTAttributers() (graph, node, edge dot.Attributer) {
	return gAttributer{}, nAttributer{}, eAttributer{}
}
