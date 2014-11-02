package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/utils"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var (
	numberOfNodes int
	numberOfEdges int
	filename      string
)

func init() {
	flag.StringVar(&filename, "filename", "path/to/generatedGraphviz.txt", "Defines how the file should be named where the generation is stored (default is ./generatedGraphviz.txt.")
	flag.IntVar(&numberOfNodes, "nodes", 2, "Defines the number of nodes that should be created. There must be more than 2 nodes.")
	flag.IntVar(&numberOfEdges, "edges", 3, "Defines the number of edges that should be used to connect the nodes. There must be more edges than nodes.")
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()
	if numberOfNodes < 2 {
		log.Fatalln("There must be more than 2 nodes.")
	}
	if numberOfEdges <= numberOfNodes {
		log.Fatalln("The number of edges must be greater than the number of nodes.")
	}
	if filename == "path/to/generatedGraphviz.txt" {
		filename = "generatedGraphviz.txt"
	}
	if err := utils.CheckIfFileExists(filename); err == nil {
		var input string
		fmt.Printf("The file \"%s\" exists. Would you like to overwrite it (y|n)?\n", filename)
		fmt.Print("\nInput: ")
		if _, err := fmt.Scanln(&input); err == nil {
			switch input {
			case "y":
				log.Println("The file gets overwritten.")
			case "n":
				log.Println("The file gets not touched.\n\nClose program.")
				os.Exit(0)
			default:
				log.Fatalln("Wrong input. Please only insert y for \"YES\" or n for \"NO\".")
			}
		} else {
			log.Fatalln("Wrong input. Please only insert y for \"YES\" or n for \"NO\".")
		}
	}
	stringBuffer := bytes.NewBufferString("graph G {\n")

	actualNumberOfEdges := 0
	randomObject := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sourceNode, destinationNode int
	for actualNumberOfEdges <= numberOfNodes {
		sourceNode = randomObject.Intn(numberOfNodes)
		sourceNode++
		destinationNode = randomObject.Intn(numberOfNodes)
		destinationNode++
		stringBuffer.WriteString(fmt.Sprintf("\t%d--%d;\n", sourceNode, destinationNode))
		actualNumberOfEdges++
	}
	stringBuffer.WriteString("}")
	//Writes the files with 0644 Unix permissions.
	ioutil.WriteFile(filename, stringBuffer.Bytes(), 0644)
	log.Println("File successfully written.")
	log.Println("Now try to generate a jpg file from the generated graphviz file.")
	cmd := exec.Command("dot", "-Tjpg", filename, "-o generatedGraphviz.jpg")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Waiting for command to finish...")
	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v", err)
		os.Exit(1)
	}
	log.Println("Image successfully created.")
}
