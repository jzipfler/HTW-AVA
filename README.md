HTW-AVA
=======

For the lecture "Architektur verteilter Anwendungen" at the htw saar.

As recommended at the GOLANG.ORG page, I created the repository on a folder
inside of a package structure (github.com/jzipfler/htw-ava).

With this structure you can download and use this package with the
"go get" command. Simply call "go get github.com/jzipfler/htw-ava" and go
will fetch the data and place them into your GOPATH folder as a GIT repo.

To run the program, use the "go build" command.
Call "go build github.com/jzipfler/htw-ava/avaStarter" to create a "avaStarter" or
maybe a "avaStarter.exe" (if you're under Windows), which can be executed from
the command line.

You can also build the generateGraphviz application with:
"go build github.com/jzipfler/htw-ava/generateGraphviz".

The example directory holds some shell scripts that executes the build
commands and start some instances with the there available NodeList and
Graphviz files.