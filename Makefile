PWD := $(shell pwd)
SRC_DIR := $(PWD)/src/go
BIN_DIR := $(PWD)/build
GOCMD := go
BUILD_CMD := go build -o $(BIN_DIR)/site
RUN_CMD := $(BIN_DIR)/site

.PHONY: build
build:
	cd $(SRC_DIR);$(BUILD_CMD)
	Build complete

.PHONY: run
run:
	$(RUN_CMD)
