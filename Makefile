.PHONY=clean

build:
	go build
clean:
	rm -rf *.o
install:
	echo "Install"
