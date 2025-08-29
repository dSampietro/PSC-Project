package main

import (
	"fmt"
	"sync"

	"golang.design/x/chann"
)

func SeqStrategy(graph Graph, max_depth int, resultCh *chann.Chann[Message]) {
	//SENTENCE GENERATION
	//resultCh := chann.New[Message]()//make(chan Message, 1000)

	// Setup channel with initial value DFS from each node
	for _, node := range graph.nodes {
		if node.label == "." { continue }
		
		msg := Message {
			//sentence: fmt.Sprintf("[FROM %s]", node.label),
			sentence: []string{fmt.Sprintf("[FROM %s]", node.label)},
			//visited: map[string]int{node.label: 1},
			depth: 0,
		}
		//node.input.In() <- msg // Start traversal with empty message

		node.GenerateSentenceSeq(msg, resultCh.In(), max_depth)
	}


	// Wait for all paths to finish: since it is seq, we wait for all GenerateSentenceSeq to finish
	//resultCh.Close()
}


func ParStrategy(graph Graph, max_depth int, resultCh *chann.Chann[Message]){
	var wg sync.WaitGroup
	
	// Setup channel with initial value DFS from each node
	for _, node := range graph.nodes {
		// we guarantee one goroutine/node => no unbounded goroutines
		node.GenerateSentence(&wg, resultCh.In(), max_depth)

		if node.label == "." { continue }
		wg.Add(1)
		
		msg := Message {
			//sentence: fmt.Sprintf("[FROM %s]", node.label),
			sentence: []string{fmt.Sprintf("[FROM %s]", node.label)},
			//visited: map[string]int{node.label: 1},
			depth: 0,
		}
		node.input.In() <- msg // Start traversal with empty message
	}


	// Wait for all paths to finish
	wg.Wait()
}