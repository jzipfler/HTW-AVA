#/bin/sh

if test ${1} = "tcp"; then
    ./fileManager -filename FileA.txt -force -ipAddress 127.1.1.1 -port 15100 -useTCP
    exit 0
fi

./fileManager -filename FileA.txt -force -ipAddress 127.1.1.1 -port 15100
