#!/bin/sh

if test ${1} = "tcp"; then
  ./avaStarter -nodeList="Nodes.txt" -isController=true -ipAddress="127.0.0.1" -port=15100 -useTCP
  exit 0
fi

./avaStarter -nodeList="Nodes.txt" -isController=true -ipAddress="127.0.0.1" -port=15100
