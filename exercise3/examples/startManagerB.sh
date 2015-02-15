#/bin/sh

if test ${1} = "tcp"; then
  ./fileManager -filename FileB.txt -force -ipAddress 127.2.2.2 -port 15102 -useTCP
  exit 0
fi

./fileManager -filename FileB.txt -force -ipAddress 127.2.2.2 -port 15102
