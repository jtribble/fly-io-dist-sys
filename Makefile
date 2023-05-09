.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help: # See: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run the simple tests for every cmd
	./scripts/make.sh test

.PHONY: mtest
mtest: ## Run the maelstrom tests for every cmd
	./scripts/make.sh mtest
