.PHONY=clean

ifeq ($(OS),Windows_NT)
	file_ext := .exe
else
	file_ext := ""
endif

executable := imgdup${file_ext}
build:
	go build -o ${executable} cmd/main/main.go
clean:
	rm -rf ${executable}
