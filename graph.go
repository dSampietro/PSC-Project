package main

import (
	"fmt"
	"sync"
	"sync/atomic"

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




type Message struct {
	sentence string
	visited map[*Node]int
	depth int
}


/*
maxInFlight is the minimum buffer size you’d need on your busiest channel to avoid ever blocking.
*/

var (
	inFlight    int64 // current # of messages sent but not yet fully handled
	maxInFlight int64 // peak value of inFlight
)

func recordPeak(n int64) {
	for {
		old := atomic.LoadInt64(&maxInFlight)
		if n <= old || atomic.CompareAndSwapInt64(&maxInFlight, old, n) {
			break
		}
	}
}

func (n *Node) GenerateSenetence(wg *sync.WaitGroup, resultCh chan<- Message, node_limit int, max_depth int) {
	go func(){
		for msg := range n.input.Out() {
			if msg.depth >= max_depth {
				atomic.AddInt64(&inFlight, -1)
				wg.Done()
                continue
            }

			if msg.visited[n] + 1 > node_limit {
				atomic.AddInt64(&inFlight, -1)
				wg.Done()
				continue
			}

			newSentence := msg.sentence + " " + n.label
			//update visited nodes of message
			newVisited := make(map[*Node]int)
			for k, v := range msg.visited {
				newVisited[k] = v
			}
			newVisited[n]++

			
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
				visitedCopy := make(map[*Node]int)
				for k, v := range newVisited {
					visitedCopy[k] = v
				}

				curr := atomic.AddInt64(&inFlight, 1)
				recordPeak(curr)

				succ.input.In() <- Message{
					sentence: newSentence,
					visited: visitedCopy,
					depth: msg.depth + 1}
			}
			wg.Done()
		}
	}()
}