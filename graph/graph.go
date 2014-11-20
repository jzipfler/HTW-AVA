package graph

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Graph struct {
	nodes int
	edges int
}

func New(numberOfNodes, numberOfEdges int) Graph {
	return Graph{numberOfNodes, numberOfEdges}
}

func (graph Graph) NumberOfNodes() int {
	return graph.nodes
}

func (graph Graph) NumberOfEdges() int {
	return graph.edges
}

func (graph *Graph) SetNumberOfNodes(numberOfNodes int) {
	// The test for -1 is done in isEdgePossible()
	graph.nodes = numberOfNodes
}

func (graph *Graph) SetNumberOfEdges(numberOfEdges int) {
	// The test for -1 is done in isEdgePossible()
	graph.edges = numberOfEdges
}

func (graph *Graph) Clear() {
	graph.nodes = 0
	graph.edges = 0
}

func (graph Graph) String() string {
	return "{ nodes: " + strconv.Itoa(graph.nodes) + ", edges: " + strconv.Itoa(graph.edges) + " }"
}

func (graph Graph) UndirectedGraph() (adjacencyMatrix [][]int, err error) {
	if graph.nodes == 0 || graph.edges == 0 {
		return nil, errors.New("Nodes and edges must not be 0.")
	}
	if graph.nodes < 0 {
		return nil, errors.New("Node graph can be generated with a negative amount of nodes.")
	}
	if graph.edges < 0 {
		return nil, errors.New("No graph can be generated with a negative amount of edges.")
	}
	actualNumberOfEdges := 0
	randomObject := rand.New(rand.NewSource(time.Now().UnixNano()))
	var destinationNode int
	adjacencyMatrix = make([][]int, graph.nodes)
	/*
		Sadly, I have to create the two dimensional array instead of doing that at runtime.
		The problem is, that I have to add values twice to the array because the graph is
		undirected. One for go there and one for going back. And because of the way back is
		not present when creating it at runtime, I have to initialize it before.
	*/
	for index := range adjacencyMatrix {
		adjacencyMatrix[index] = make([]int, graph.nodes)
	}
	for index := range adjacencyMatrix {
		if index == 0 {
			continue
		}
		destinationNode = randomObject.Intn(index)
		if _, err := graph.isEdgePossible(index, destinationNode, adjacencyMatrix); err != nil {
			//log.Println(err)
			continue
		}
		if err := graph.addUndirectedEdge(index, destinationNode, adjacencyMatrix); err != nil {
			log.Fatalln(err)
		}
		actualNumberOfEdges++
	}
	for actualNumberOfEdges < graph.edges {
		sourceNode := randomObject.Intn(graph.nodes)
		for destinationNode := range adjacencyMatrix[sourceNode] {
			if possible, _ := graph.isEdgePossible(sourceNode, destinationNode, adjacencyMatrix); possible {
				if err := graph.addUndirectedEdge(sourceNode, destinationNode, adjacencyMatrix); err != nil {
					log.Fatalln(err)
				}
				actualNumberOfEdges++
			}
		}
	}
	return
}

func (graph Graph) UndirectedGraphAsDotLanguageString() (string, error) {
	stringBuffer := bytes.NewBufferString("graph G {\n")
	adjacencyMatrix, err := graph.UndirectedGraph()
	if err != nil {
		return "", err
	}
	for fromIndex := range adjacencyMatrix {
		for toIndex := range adjacencyMatrix[fromIndex] {
			if adjacencyMatrix[fromIndex][toIndex] == 1 {
				from := fromIndex + 1
				to := toIndex + 1
				stringBuffer.WriteString(fmt.Sprintf("\t%d -- %d;\n", from, to))
				//fmt.Printf("FROM: %d TO: %d\n", fromIndex, toIndex)
				adjacencyMatrix[fromIndex][toIndex] = 0
				adjacencyMatrix[toIndex][fromIndex] = 0
			}
		}
	}
	stringBuffer.WriteString("}")
	return stringBuffer.String(), nil
}

func (graph Graph) DirectedGraph() (adjacencyMatrix [][]int, err error) {
	if graph.nodes == 0 || graph.edges == 0 {
		return nil, errors.New("Nodes and edges must not be 0.")
	}
	if graph.nodes < 0 {
		return nil, errors.New("Node graph can be generated with a negative amount of nodes.")
	}
	if graph.edges < 0 {
		return nil, errors.New("No graph can be generated with a negative amount of edges.")
	}
	actualNumberOfEdges := 0
	randomObject := rand.New(rand.NewSource(time.Now().UnixNano()))
	var destinationNode int
	adjacencyMatrix = make([][]int, graph.nodes)
	/*
		Sadly, I have to create the two dimensional array instead of doing that at runtime.
		The problem is, that I have to add values twice to the array because the graph is
		undirected. One for go there and one for going back. And because of the way back is
		not present when creating it at runtime, I have to initialize it before.
	*/
	for index := range adjacencyMatrix {
		adjacencyMatrix[index] = make([]int, graph.nodes)
	}
	for index := range adjacencyMatrix {
		if index == 0 {
			continue
		}
		destinationNode = randomObject.Intn(index)
		if _, err := graph.isEdgePossible(index, destinationNode, adjacencyMatrix); err != nil {
			//log.Println(err)
			continue
		}
		if err := graph.addDirectedEdge(index, destinationNode, adjacencyMatrix); err != nil {
			log.Fatalln(err)
		}
		actualNumberOfEdges++
	}
	for actualNumberOfEdges <= graph.edges {
		sourceNode := randomObject.Intn(graph.nodes)
		for destinationNode := range adjacencyMatrix[sourceNode] {
			if possible, _ := graph.isEdgePossible(sourceNode, destinationNode, adjacencyMatrix); possible {
				if err := graph.addDirectedEdge(sourceNode, destinationNode, adjacencyMatrix); err != nil {
					log.Fatalln(err)
				}
				actualNumberOfEdges++
			}
		}
	}
	return
}

func (graph Graph) DirectedGraphAsDotLanguageString() (string, error) {
	adjacencyMatrix, err := graph.DirectedGraph()
	if err != nil {
		return "", err
	}
	stringBuffer := bytes.NewBufferString("digraph G {\n")
	for fromIndex := range adjacencyMatrix {
		for toIndex := range adjacencyMatrix[fromIndex] {
			if adjacencyMatrix[fromIndex][toIndex] == 1 {
				from := fromIndex + 1
				to := toIndex + 1
				stringBuffer.WriteString(fmt.Sprintf("\t%d -> %d;\n", from, to))
				//fmt.Printf("FROM: %d TO: %d\n", fromIndex, toIndex)
				adjacencyMatrix[fromIndex][toIndex] = 0
			}
		}
	}
	stringBuffer.WriteString("}")
	return stringBuffer.String(), nil
}

func (graph *Graph) addUndirectedEdge(source, destination int, adjacencyMatrix [][]int) error {
	if _, err := graph.isEdgePossible(source, destination, adjacencyMatrix); err != nil {
		return err
	}
	adjacencyMatrix[source][destination] = 1
	adjacencyMatrix[destination][source] = 1
	return nil
}

func (graph *Graph) addDirectedEdge(source, destination int, adjacencyMatrix [][]int) error {
	if _, err := graph.isEdgePossible(source, destination, adjacencyMatrix); err != nil {
		return err
	}
	adjacencyMatrix[source][destination] = 1
	return nil
}

func (graph Graph) isEdgePossible(source, destination int, adjacencyMatrix [][]int) (bool, error) {
	if source < 0 {
		return false, errors.New("Source may not be negative.")
	}
	if destination < 0 {
		return false, errors.New("Destination my not be negative.")
	}
	if source > len(adjacencyMatrix) {
		return false, errors.New("Source may not be larger than the array.")
	}
	if destination > len(adjacencyMatrix) {
		return false, errors.New("Destination may not be larger than the array.")
	}
	if source == destination {
		return false, errors.New("Source and destination have to be different.")
	}
	for index := range adjacencyMatrix {
		if len(adjacencyMatrix[index]) < source {
			return false, errors.New("Source may not be larger than the second dimension of the array.")
		}
		if len(adjacencyMatrix[index]) < destination {
			return false, errors.New("Destination may not be larger than the second dimension of the array.")
		}
	}
	if adjacencyMatrix[source][destination] == 1 {
		return false, errors.New(fmt.Sprintf("A edge already exists for \"%d -- %d\"", source, destination))
	}
	return true, nil
}
