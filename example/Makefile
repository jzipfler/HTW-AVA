.PHONY: all, clean

all: avaStarter generateGraphviz

avaStarter:
	go build -v github.com/jzipfler/htw-ava/avaStarter

generateGraphviz:
	go build -v github.com/jzipfler/htw-ava/generateGraphviz

clear: clean

clean:
	@echo "Another test"
	/bin/rm -f avaStarter
	/bin/rm -f generateGraphviz
	/bin/rm -f nodeOutput*.txt
	/bin/rm -f *_belives.txt
