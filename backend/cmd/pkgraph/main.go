package main

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/kr/pretty"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/script"
)

func main() {
	script.Run(fx.Invoke(run))
}

func run(d fx.DotGraph) error {
	// fmt.Print(d)

	graphast, err := gographviz.Parse([]byte(d))
	if err != nil {
		return err
	}

	fxgraph := gographviz.NewGraph()
	graph := gographviz.NewGraph()

	graph.AddSubGraph("", "cluster_resources", nil)
	graph.AddSubGraph("", "cluster_services", nil)
	graph.AddSubGraph("", "cluster_infrastructure", nil)

	// graph.Nodes = fxgraph.Nodes
	// graph.Relations.ChildToParents = fxgraph.Relations.ChildToParents

	gographviz.Analyse(graphast, fxgraph)

	flat := fxgraph.SubGraphs.Sorted()

	// pretty.Println(flat)
	// pretty.Println(fxgraph.Nodes.Lookup)

	convattr := func(a gographviz.Attrs) map[string]string {
		b := map[string]string{}
		for k, v := range a {
			b[string(k)] = v
		}
		return b
	}

	for _, v := range flat {
		pkg := v.Attrs["label"]

		if strings.Contains(pkg, "services") {
			if err := graph.AddSubGraph("cluster_services", v.Name, convattr(v.Attrs)); err != nil {
				return err
			}

			pretty.Println(pkg, fxgraph.Relations.ParentToChildren[v.Name])

			for k := range fxgraph.Relations.ParentToChildren[v.Name] {
				n := fxgraph.Nodes.Lookup[k]

				fmt.Println("EDGE", v.Name, n.Name)

				// pretty.Println(v.Name, "CHILD", fxgraph.Nodes.Lookup[k])
				graph.AddEdge(v.Name, n.Name, true, convattr(n.Attrs))
			}
		}
	}

	graph.Directed = true

	g := graph.String()

	fmt.Println(g)

	return nil
}
