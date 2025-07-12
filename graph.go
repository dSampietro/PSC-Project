package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"golang.design/x/chann"
)

type Node struct {
	label 		string
	input 		*chann.Chann[Message]
	successors 	[]*Node
}

func NewNode(s string) *Node {
	return &Node{
		label: s,
		input: chann.New[Message](), 	//add buffer to make async
		successors: []*Node{},
	}
}

type Graph struct {
	nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node, 100),
	}
}


func (g *Graph) AddNode(label string) {
	g.nodes[label] = NewNode(label)
}


func (g *Graph) AddEdge(fromLabel, toLabel string) {
	fromNode := g.nodes[fromLabel]
	toNode := g.nodes[toLabel]
	if fromNode != nil && toNode != nil {	//prevent from adding same edge multiple times
		if !slices.Contains(fromNode.successors, toNode){
			fromNode.successors = append(fromNode.successors, toNode)
		}
	}
}


func (g *Graph) ToDot() string {
	var builder strings.Builder
	
	builder.WriteString("digraph {\n")
	for _, node := range g.nodes {
		if node.successors != nil {
			for _, succ  := range node.successors {
				line := fmt.Sprintf("\t\"%s\" -> \"%s\"\n", node.label, succ.label)
				builder.WriteString(line)
			}
		}
	}
	builder.WriteString("}\n")

	return builder.String()
}




type Message struct {
	sentence string
	visited map[string]int
	depth int
}


func (n *Node) GenerateSenetence(wg *sync.WaitGroup, resultCh chan<- Message, node_limit int, max_depth int) {
	go func(){
		for msg := range n.input.Out() {
			if msg.depth >= max_depth {
				wg.Done()
                continue
            }

			if msg.visited[n.label] + 1 > node_limit {
				wg.Done()
				continue
			}

			newSentence := msg.sentence + " " + n.label
			//update visited nodes of message
			newVisited := make(map[string]int, len(msg.visited))
			for k, v := range msg.visited {
				newVisited[k] = v
			}
			newVisited[n.label]++

			
			if len(n.successors) == 0 {	//terminal node
				resultCh <- Message{
					sentence: newSentence,
					visited: newVisited,
					depth: msg.depth + 1}
				wg.Done()
				continue
			} 

			//forward to successors
			for _, succ := range n.successors {
				wg.Add(1)

				//clone visited list for each successor, to avoid mutable sharing
				visitedCopy := make(map[string]int, len(msg.visited))
				for k, v := range newVisited {
					visitedCopy[k] = v
				}

				succ.input.In() <- Message{
					sentence: newSentence,
					visited: visitedCopy,
					depth: msg.depth + 1}
			}
			wg.Done()
		}
	}()
}