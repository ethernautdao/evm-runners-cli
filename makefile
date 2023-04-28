BINARY_NAME=evm-runners
INSTALL_DIR=${HOME}/.evm-runners

all: clean build

build:
	go build -o ${BINARY_NAME} main.go

install:
	mkdir -p ${INSTALL_DIR}
	cp ${BINARY_NAME} ${INSTALL_DIR}

clean:
	go clean
	rm -f ${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME}
