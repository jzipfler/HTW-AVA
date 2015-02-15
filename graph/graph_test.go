package graph_test

import (
	"github.com/jzipfler/htw-ava/graph"
	"testing"
)

func TestGetterAndSetter(t *testing.T) {
	graphObject := graph.New(0,0)
	
	if graphObject.NumberOfEdges() != 0 || graphObject.NumberOfNodes() != 0 {
		t.Error("The initiation of the object failed. Both values must be zero.")
	}	
	graphObject.SetNumberOfNodes(10)
	graphObject.SetNumberOfEdges(10)
	if graphObject.NumberOfEdges() != 10 || graphObject.NumberOfNodes() != 10 {
		t.Errorf("The value of nodes and edges have to be 10 but are nodes(%d), edges(%d).",graphObject.NumberOfEdges(), graphObject.NumberOfNodes())
	}
}

func TestClearingFunction(t *testing.T) {
	graphObject := graph.New(1,1)
	graphObject.Clear()
	nodes := graphObject.NumberOfNodes()
	edges := graphObject.NumberOfEdges()
	
	if nodes != 0 || edges != 0 {
		t.Errorf("Nodes (%d) or edges (%d) are not 0, so clear does not worked.", nodes, edges)
	}
}

func TestEmptyObjectBehaviour(t *testing.T) {
	graphObject := graph.New(0,0)
	if _, err := graphObject.UndirectedGraph(); err == nil {
		t.Error("No error occured on the UndirectedGraph function with an empty graph object.")
	}
	if _, err := graphObject.UndirectedGraphAsDotLanguageString(); err == nil {
		t.Error("No error occured on the UndirectedGraphAsDotLanguageString function with an empty graph object.")
	}
	if _, err := graphObject.DirectedGraph(); err == nil {
		t.Error("No error occured on the DirectedGraph function with an empty graph object.")
	}	
	if _, err := graphObject.DirectedGraphAsDotLanguageString(); err == nil {
		t.Error("No error occured on the DirectedGraphAsDotLanguageString function with an empty graph object.")
	}	
}

func TestNegativeValues(t *testing.T) {
	graphObject := graph.New(-1,-1)
	if _, err := graphObject.DirectedGraph(); err == nil {
		t.Error("The nodes and edges was set to -1 and no error occured for DirectedGraph.")
	}
	if _, err := graphObject.UndirectedGraph(); err == nil {
		t.Error("The nodes and edges was set to -1 and no error occured for UndirectedGraph.")
	}
}