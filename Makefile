include Makefile.in

.PHONY: all clean test run build upgrade help

all: 			# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean:			# clean-up environment
	rm -rf $(BIN)

test:			# run test
	$(GOFMT) -w -s $(SRC)
	$(GO) test -cover -failfast -timeout 2s ./...

run: $(BIN)		# run in the local environment
	./$(BIN) -Fpvvv

build: $(BIN)	# build the binary/library

upgrade:		# upgrade all the necessary packages
	pre-commit autoupdate
	$(GO) get -u all

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(BIN): test

%: %.go
	$(GO) mod tidy
	$(GO) build -ldflags="-s -w" -o $@ $<

.PHONY: stop image log
log:			# show the logs inside docker
	docker-compose logs -f

stop:			# stop the docker running
	docker-compose stop

image:			# build the docker image
	docker-compose stop
	docker-compose up --build -d
