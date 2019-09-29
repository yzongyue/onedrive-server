BINARY=onedrive-server
test:
	go test -v -cover
build:
	go build -o ${BINARY}
clean:
	if [ -a ${BINARY} ] ; then rm ${BINARY}; fi
