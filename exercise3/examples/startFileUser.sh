#/bin/sh

if test ${1} = "tcp"; then
  TCP=true
  shift
else
  TCP=false
fi

if test ${#} -ne 1; then
	./fileUser -ipAddress 127.5.5.5 -port 15200 -managerA 127.1.1.1:15100 -managerB 127.2.2.2:15102 -useTCP=${TCP}
else
	./fileUser -ipAddress 127.6.6.6 -port 15202 -managerA 127.1.1.1:15100 -managerB 127.2.2.2:15102 -useTCP=${TCP}
fi
