BINARY_NAME=evm-runners
ALT_NAME=evmr
INSTALL_DIR=${HOME}/.evm-runners

all: clean build

build: 
	go build -o ${BINARY_NAME} main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64 main.go

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}-linux-arm64 main.go
	
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-darwin-amd64 main.go

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-darwin-arm64 main.go

prod: build-linux build-linux-arm64 build-macos build-macos-arm64

install:
	mkdir -p ${INSTALL_DIR}
	cp ${BINARY_NAME} ${INSTALL_DIR}
	make symlink

symlink:
	ln -s ${INSTALL_DIR}/${BINARY_NAME} ${INSTALL_DIR}/${ALT_NAME}

clean:
	go clean
	rm -f ${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME} ${INSTALL_DIR}/${ALT_NAME}
