package main

import (
	"fmt"
	"sync"
)

type Node struct {
	label 		string
	input 		chan string
	successors 	[]*Node
}

func NewNode(s string) *Node {
	return &Node{
		label: s,
		input: make(chan string, 100), 	//add buffer to make async
		successors: []*Node{},
	}
}


type Message struct {
	sentence string
	visited map[*Node]int
}

func (n *Node) GenerateSenetence(wg *sync.WaitGroup, resultCh chan<- string) {
	go func(){
		for msg := range n.input {
			newMsg := msg + " " + n.label
			
			if len(n.successors) == 0 {	//terminal node
				resultCh <- newMsg
				wg.Done()
			} else {
				for _, succ := range n.successors {
					wg.Add(1)
					succ.input <- newMsg
				}
				wg.Done()
			}
		}
	}()
}


type Graph struct {
	nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
	}
}


func (g *Graph) AddNode(label string) *Node {
	node := &Node{
		label:      label,
		input:      make(chan string),
		successors: []*Node{},
	}
	g.nodes[label] = node
	return node
}


func (g *Graph) AddEdge(fromLabel, toLabel string) {
	fromNode := g.nodes[fromLabel]
	toNode := g.nodes[toLabel]
	if fromNode != nil && toNode != nil {
		fromNode.successors = append(fromNode.successors, toNode)
	}
}


func (g *Graph) PrettyPrint(){
	fmt.Println("digraph {")
	for _, node := range g.nodes {
		if node.successors != nil {
			for _, succ  := range node.successors {
				fmt.Printf("\t%s -> %s\n", node.label, succ.label)
			}
		}
	}
	fmt.Printf("}\n")
}