BINARY_NAME=evm-runners
INSTALL_DIR=/usr/local/bin/

all: clean build

build:
	go build -o ${BINARY_NAME} main.go

install:
	cp ${BINARY_NAME} ${INSTALL_DIR}

clean:
	go clean
	rm -f ${BINARY_NAME}
