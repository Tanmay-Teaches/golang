package main

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"
)

const INFINITY = ^uint(0)

type Node struct {
	Name  string
	links []Edge
}

type Edge struct {
	from *Node
	to   *Node
	cost uint
}
type Graph struct {
	nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{nodes: map[string]*Node{}}
}
func (g *Graph) AddNodes(names ...string) {
	for _, name := range names {
		if _, ok := g.nodes[name]; !ok {
			g.nodes[name] = &Node{Name: name, links: []Edge{}}
		}
	}
}

//Assume all links are undirected
func (g *Graph) AddLink(a, b string, cost int) {
	aNode := g.nodes[a]
	bNode := g.nodes[b]
	//Link a to b
	aNode.links = append(aNode.links, Edge{from: aNode, to: bNode, cost: uint(cost)})
	//Link b to a
	bNode.links = append(aNode.links, Edge{from: bNode, to: aNode, cost: uint(cost)})

}
func (g *Graph) Dijkstra(source string) (map[string]uint, map[string]string) {
	dist, prev := map[string]uint{}, map[string]string{}

	for _, node := range g.nodes {
		dist[node.Name] = INFINITY
		prev[node.Name] = ""
	}
	visited := map[string]bool{}
	dist[source] = 0
	for u := source; u != ""; u = getClosestNonVisitedNode(dist, visited) {
		uDist := dist[u]
		for _, link := range g.nodes[u].links {
			if _, ok := visited[link.to.Name]; ok {
				continue
			}
			alt := uDist + link.cost
			v := link.to.Name
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
			}
		}
		visited[u] = true
	}
	return dist, prev
}

func getClosestNonVisitedNode(dist map[string]uint, visited map[string]bool) string {
	lowestCost := INFINITY
	lowestNode := ""
	for key, dis := range dist {
		if _, ok := visited[key]; dis == INFINITY || ok {
			continue
		}
		if dis < lowestCost {
			lowestCost = dis
			lowestNode = key
		}
	}
	return lowestNode
}

func main() {
	g := NewGraph()
	g.AddNodes("a", "b", "c", "d", "e")
	g.AddLink("a", "b", 6)
	g.AddLink("a", "d", 1)
	g.AddLink("d", "b", 2)
	g.AddLink("d", "e", 1)
	g.AddLink("e", "b", 2)
	g.AddLink("e", "c", 5)
	g.AddLink("c", "b", 5)
	dist, prev := g.Dijkstra("a")
	fmt.Println(DijkstraString(dist, prev))
}

func DijkstraString(dist map[string]uint, prev map[string]string) string {
	buf := &bytes.Buffer{}
	writer := tabwriter.NewWriter(buf, 1, 5, 2, ' ',0)
	writer.Write([]byte("Node\tDistance\tPrevious Node\t\n"))
	for key, value := range dist {
		writer.Write([]byte(key + "\t"))
		writer.Write([]byte(strconv.FormatUint(uint64(value), 10) + "\t"))
		writer.Write([]byte(prev[key] + "\t\n"))
	}
	writer.Flush()
	return buf.String()
}
