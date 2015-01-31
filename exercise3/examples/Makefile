.PHONY: all clean clear fileManager fileUser

all: fileManager fileUser

fileManager:
	go build -v github.com/jzipfler/htw-ava/exercise3/fileManager

fileUser:
	go build -v github.com/jzipfler/htw-ava/exercise3/fileUser

clear: clean

clean:
	@echo "Remove generated files if available"
	/bin/rm -f fileManager
	/bin/rm -f fileUser

