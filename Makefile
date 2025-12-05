.PHONY: build test pr

SUBMODULES = $(wildcard */)

build:
	bash ./container/run.sh multimod build

test:
	bash ./container/run.sh multimod test

lint:
	bash ./container/run.sh multimod lint

deps:
	bash ./container/run.sh multimod update

pr:
	bash ./container/run.sh multimod annotate
