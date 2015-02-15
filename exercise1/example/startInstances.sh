#!/bin/sh

if ${1} = "tcp"; then
  # Do not use the 8th one because the Nodes.txt holds a error on this id.
  for i in 1 2 3 4 5; do
    ./avaStarter -nodeList="Nodes.txt" -graphvizFile="Graphviz.txt" -isController=false -useTCP -id ${i} > nodeOutput${i}.txt &
  done
  exit 0
fi

for i in 1 2 3 4 5; do
    ./avaStarter -nodeList="Nodes.txt" -graphvizFile="Graphviz.txt" -isController=false -id ${i} > nodeOutput${i}.txt &
done
