.DEFAULT_GOAL := build
EXECUTABLE = book
GREEN = "\033[1;38;5;048m"
RED = "\033[1;38;5;196m"
CLS = "\033[0m"

# test:
# 	@echo $(GREEN)"RUNNING TESTS... 🤓"$(CLS)
# 	@if go test -tags skip ./...; then \
# 		echo $(GREEN)"ALL TESTS PASSED 👍"$(CLS); \
# 	else \
# 		echo $(RED)"CHECK FAILING TEST(S) 👎"$(CLS); \
# 		exit 1; \
# 	fi
# .PHONY:test

fmt:
	@echo $(GREEN)"RUNNING FORMATTER... 🐸"$(CLS)
	@go fmt ./...
.PHONY:fmt

lint: fmt
	@echo $(GREEN)"RUNNING LINTER... 🦁"$(CLS)
	@golint ./...
.PHONY:lint

vet: lint
	@echo $(GREEN)"RUNNING VETTER... 🐹"$(CLS)
	@go vet ./...
	@shadow ./...
.PHONY:vet

build: vet
	@go build && \
	echo $(GREEN)"BUILD SUCCEEDED 👌"$(CLS)
.PHONY:build

clean:
	@if [[ -a $(EXECUTABLE) ]]; then \
		rm $(EXECUTABLE) && echo $(GREEN)"CLEANED 👍"$(CLS); \
	else \
		 echo $(RED)"NOTHING TO CLEAN 👎"$(CLS); \
	fi
.PHONY:clean
