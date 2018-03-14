package grapher

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/thomaspurchas/monzo-test/crawler"
)

func BuildGraph(uchan <-chan *crawler.URLContext) graph.Graph {
	nodes := make(map[string]Node)
	g := newGraph()

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
