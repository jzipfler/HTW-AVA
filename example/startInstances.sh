#!/bin/sh

# Do not use the 8th one because the Nodes.txt holds a error on this id.
for i in 1 2 3 4 5 6 7 9; do
	./avaStarter -nodeList="Nodes.txt" -isController=false -id ${i} > nodeOutput${i}.txt &
done
