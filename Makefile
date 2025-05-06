# Makefile

APP_NAME = cookrobot
GOOS ?= darwin
GOARCH ?= arm64

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME) *.go
