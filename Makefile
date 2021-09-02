default: build

build: ## build slim-orderbook
		'$(CURDIR)/scripts/src.build.sh'
clean:
		'$(CURDIR)/scripts/src.cleanup.sh'

.PHONY: default help build clean