package main

import (
	"fmt"
	"sync"

	"golang.design/x/chann"
)

func SeqStrategy(graph Graph, max_depth int, resultCh *chann.Chann[Message]) {
	// Setup channel with initial value DFS from each node
	for _, node := range graph.nodes {
		if node.label == "." { continue }
		
		msg := Message {
			sentence: []string{fmt.Sprintf("[FROM %s]", node.label)},
			depth: 0,
		}
		//node.input.In() <- msg // Start traversal with empty message

		node.GenerateSentenceSeq(msg, resultCh.In(), max_depth)
	}

	// Wait for all paths to finish: since it is seq, we wait for all GenerateSentenceSeq to finish
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
			sentence: []string{fmt.Sprintf("[FROM %s]", node.label)},
			depth: 0,
		}
		node.input.In() <- msg // Start traversal with empty message
	}


	// Wait for all paths to finish
	wg.Wait()
}