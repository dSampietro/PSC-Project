package main

import (
	"fmt"
	"sync"
)

type Node struct {
	label 		string
	input 		chan Message
	successors 	[]*Node
}

func NewNode(s string) *Node {
	return &Node{
		label: s,
		input: make(chan Message, 100), 	//add buffer to make async
		successors: []*Node{},
	}
}


type Message struct {
	sentence string
	visited map[*Node]int
}


func NewMessage(s string) Message {
	return Message {
		sentence: s,
		visited: make(map[*Node]int),
	}
}


//TODO: abort when > N
func (n *Node) GenerateSenetence(wg *sync.WaitGroup, resultCh chan<- Message) {
	go func(){
		for msg := range n.input {

			//update visited nodes of message
			newVisited := make(map[*Node]int)
			for k, v := range msg.visited {
				newVisited[k] = v
			}
			newVisited[n]++

			newMsg := Message{
				sentence: msg.sentence + " " + n.label,
				visited: newVisited,
			}
			
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


func (g *Graph) AddNode(label string) {
	g.nodes[label] = NewNode(label)
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