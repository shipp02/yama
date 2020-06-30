PWD := $(shell pwd)

SRC_DIR := $(PWD)/src/go
BIN_DIR := $(PWD)/build
TEST_DIR := $(PWD)/src/test

GOCMD := go
BUILD_CMD := go build -o $(BIN_DIR)/site
RUN_CMD := $(BIN_DIR)/site

.PHONY: build
build:
	cd $(SRC_DIR);$(BUILD_CMD)
	echo "Build complete"

.PHONY: run
run:
	make build
	$(RUN_CMD)

.PHONY: doc
doc:
	go doc $(SRC_DIR)
