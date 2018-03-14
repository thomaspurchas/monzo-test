# monzo.com crawler

## Setup

Setup can be done with `go get ./...` in this directory, or even `go get github.com/thomaspurchas/monzo-test`.

For best result please also have [GraphViz][1] installed (with `dot` in your `PATH`), on OSX it is
available from HomeBrew (`brew install graphviz`). For other platforms visit the
[download page](https://www.graphviz.org/download/).

## Running

Once setup with dependences in your `GOPATH` run `go run monzo.go`. This will output log data to
`stderr` and dot output to `stdout`.

To render visualisation (you need [GraphViz][1]) run `go run monzo.go | dot -Tpdf output.pdf`. You
can render to other formats, but PDF is advised due to the size of the graph.

[1]:https://www.graphviz.org/
