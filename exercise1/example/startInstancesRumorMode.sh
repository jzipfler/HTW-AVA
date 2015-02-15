#!/bin/sh

# Do not use the 8th one because the Nodes.txt holds a error on this id.
for i in 1 2 3 4 5; do
	./avaStarter -nodeList="Nodes.txt" -graphvizFile="Graphviz.txt" -rumorExperiment -id ${i} > nodeOutput${i}.txt &
done
