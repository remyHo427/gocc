MYCC_SRCPATH=./cmd/cc
BIN_PATH=./bin
TEST_DIR=./tests

.DEFAULT_GOAL=make
.PHONY: clean, test

make:
	@go run ./scripts/main.go build mycc

test:
	@go run ./scripts/main.go test mycc
	@(cp $(BIN_PATH)/mycc $(TEST_DIR))
	@(cd $(TEST_DIR) && ./test_all.sh)
