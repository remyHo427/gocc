MYCC_SRCPATH=./cmd/cc
BIN_PATH=./bin
TEST_DIR=./tests

.DEFAULT_GOAL=make
.PHONY: test

make:
	@go run ./scripts/main.go build mycc

test:
	@go test ./cmd/... | grep -E "^[^?]"

compile_test: $(BIN_PATH)/mycc
	@(cp $(BIN_PATH)/mycc $(TEST_DIR))
	@(cd $(TEST_DIR) && ./test_all.sh)

