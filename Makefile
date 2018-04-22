FILENAME=drone-control

build:
	go build -o ${FILENAME} main.go

build_pi:
	GOOS=linux GOARCH=arm GOARM=6 go build -o ${FILENAME} main.go

mocks:
	go generate ./...
