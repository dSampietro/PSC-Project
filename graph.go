package main

import "sync"

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


type Graph map[int]*Node

func (g Graph) AddNode(label string) {
	
}
func (g Graph) AddEdge() {}