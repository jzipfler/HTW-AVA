#!/bin/sh

# For exercise 2, the number of customer and company nodes must be variable.
# So I have to introduce comand line parameters for this script.

APPLICATION_NAME=$0

# The help function...
printHelp() {
  echo "Usage of $(basename ${APPLICATION_NAME}):
    -h    Prints this help message.
    -c    Sets the number of companies.
    -C    Sets the number of customers.
    -n    Sets the name of the node list."
}

# Define a default number of customer and companies.
COMPANY=2
CUSTOMER=3

VERBOSE=0
NODE_LIST="Nodes.txt"

# Variable that is used from getopts to parse the set number of argument.
# Has to be reseted when used twice.
OPTIND=1

# Big C for the customer --> The customer is the king. ;-)
while getopts "h?vc:C:n:" PARAMETERS; do
  case "$PARAMETERS" in 
    h|\?)
      printHelp
      exit 0
      ;;
    v)
      VERBOSE=1
      ;;
    b)
      COMPANY=${OPTARG}
      ;;
    C)
      CUSTOMER=${OPTARG}
      ;;
    n)
      NODE_LIST=${OPTARG}
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

exit 0

# Do not use the 8th one because the Nodes.txt holds a error on this id.
for i in 1 2 3 4 5; do
	./avaStarter -nodeList="Nodes.txt" -graphvizFile="Graphviz.txt" -isController=false -id ${i} > nodeOutput${i}.txt &
done
