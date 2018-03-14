package grapher

import (
	"fmt"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/thomaspurchas/monzo-test/crawler"
)

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

type Graph struct {
	g graph.Graph
}

func BuildGraph(uchan <-chan *crawler.URLContext) *simple.DirectedGraph {
	nodes := make(map[string]Node)
	g := simple.NewDirectedGraph(1.0, 1.0)

	for u := range uchan {
		n := Node{id: g.NewNodeID(), name: u.NormalisedURL.String()}

		if _, exists := nodes[u.NormalisedURL.String()]; !exists {
			nodes[u.NormalisedURL.String()] = n
			g.AddNode(n)
		} else {
			n = nodes[u.NormalisedURL.String()]
		}

		if u.NormalisedSourceURL != nil {
			if p, exists := nodes[u.NormalisedSourceURL.String()]; exists {
				e := simple.Edge{F: p, T: n}
				if p.ID() != n.ID() {
					g.SetEdge(e)
				}
			}
		}
	}
	return g
}
