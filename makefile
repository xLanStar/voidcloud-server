.PHONY: all

all: build exec

build:
	go build -o bin

exec:
	cd ./bin; ./voidcloud-server.exe