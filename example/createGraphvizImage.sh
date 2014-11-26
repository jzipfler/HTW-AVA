#!/bin/bash
#
# This script is used to convert a graphviz definition (dot definition)
# to eigther jpg, png, pdf or svg.
#

NAME=${0##*/}

# Check if the dot program is installed
# if not, we can directly exit
which dot > /dev/null 2>&1
if test $? -ne 0; then
	echo -e "The dot tool seems not to be installed or is not in PATH."
	exit 1
fi


if test $# -ne 1; then
	echo -e "The script \"${NAME}\" needs one parameter."
	echo -e "\t--> The graphviz (dot) file."
	exit 2
fi
if test "${1}" = "--help"; then
	echo -e "The script \"${NAME}\" needs one parameter."
	echo -e "\t--> The graphviz (dot) file."
	exit 2
fi

if ! test -e ${1} || ! test -f ${1}; then
	echo -e "\n${1} existiert nicht oder ist keine regulaere Datei."
	exit 3
fi

EXIT="Abbruch/Beenden"
PS3="
${1} als Eingabe gelesen.
Wählen Sie das Ausgabeformat:"

echo "Verfügbare Ausgabeformate:"
select AUSWAHL in jpg png pdf svg ${EXIT}; do
	if test "${AUSWAHL}" = "${EXIT}"; then
		echo "Programm wird beendet"
		break
	fi
	OUTPUT_FILE="${1%.*}.${AUSWAHL}"
	dot -T${AUSWAHL} ${1} -o ${OUTPUT_FILE}
	echo -e "Datei geschrieben nach: ${OUTPUT_FILE}"
	break
done
exit 0
