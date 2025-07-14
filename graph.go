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
	sentence []string
	depth int
}


func (n *Node) GenerateSentence(wg *sync.WaitGroup, resultCh chan<- Message, max_depth int) {
	go func(){
		for msg := range n.input.Out() {
			//abort message if longer than max_depth
			if msg.depth >= max_depth {
				wg.Done()
                continue
            }

			//newSentence := msg.sentence + " " + n.label
			newSentence := append(msg.sentence, n.label)

			
			if len(n.successors) == 0 {	//terminal node
				clonedSentence := make([]string, len(newSentence))
				copy(clonedSentence, newSentence)

				resultCh <- Message{
					sentence: clonedSentence,
					depth: msg.depth + 1}
				wg.Done()
				continue
			} 

			//forward to successors
			for _, succ := range n.successors {
				wg.Add(1)
				clonedSentence := make([]string, len(newSentence))
				copy(clonedSentence, newSentence)

				succ.input.In() <- Message{
					sentence: clonedSentence,
					depth: msg.depth + 1}
			}
			wg.Done()
		}
	}()
}