.PHONY=clean

build:
	go build -o main cmd/main/main.go
clean:
	rm -rf main
