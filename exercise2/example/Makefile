.PHONY: all, clean, clear, avaStarter2 , generateGraphviz

all: avaStarter2 generateGraphviz

avaStarter2:
	go build -v github.com/jzipfler/htw-ava/avaStarter2

generateGraphviz:
	go build -v github.com/jzipfler/htw-ava/generateGraphviz

clear: clean

clean:
	@echo "Another test"
	/bin/rm -f avaStarter2
	/bin/rm -f generateGraphviz
	/bin/rm -f nodeOutput*.txt
