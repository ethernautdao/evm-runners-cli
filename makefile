BINARY_NAME=evmrunners

all: clean build

build:
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm -f ${BINARY_NAME}