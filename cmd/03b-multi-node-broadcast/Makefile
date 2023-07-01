.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help: # See: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: bin
bin: ## Build the Go binary
	go build -o bin/node

.PHONY: run
run: ## Run main.go
	go run .

.PHONY: test
test: ## Run the simple test with test.json
	jq -c '.[]' < test.json | go run .

.PHONY: mtest
mtest: bin ## Run the maelstrom test
	rm -rf store
	../../.maelstrom/maelstrom test -w broadcast --bin bin/node --node-count 5 --time-limit 20 --rate 10
