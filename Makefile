NAME = gowork
INSTALL_PATH ?= $(GOPATH)/bin

build:
	go build -o bin/$(NAME)

install: build
	cp bin/$(NAME) $(INSTALL_PATH)/