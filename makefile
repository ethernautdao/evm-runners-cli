BINARY_NAME=evm-runners

all: clean build

build:
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm -f ${BINARY_NAME}