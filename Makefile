.PHONY: $(HYLOBLOG) $(REPOSITORIES) $(DB) $(BIN) 

DOCKER = $(SUDO) docker
GO = go
BIN = ${CURDIR}/bin
SOURCES := $(shell find $(CURDIR) -name '*.go')
HYLOBLOG = $(BIN)/hyloblog

$(HYLOBLOG): $(BIN) $(SOURCES) db get build.sh
	@printf 'BUILD\t$@\n'
	@./build.sh $@ $(BUILDARGS)

get: go.mod go.sum
	@printf 'GO\tmod tidy\n'
	@go mod tidy

test: get
	@printf 'GO\ttest\n'
	@go test ./... -v

$(BIN):
	@mkdir -p $@

DBDIR = internal/model
dbfiles := $(shell find $(DBDIR) -name '*.sql')
db: $(DBDIR)/sqlc.yaml $(dbfiles)
	@printf 'SQLC\t$<\n'
	@sqlc generate -f $<

up: $(HYLOBLOG) test
	@echo 'launching docker containers...'
	$(DOCKER) compose up --build

github:
	smee -u "https://smee.io/D9yWYTiYzjBhfU3O" -t "http://localhost:7999/gh/installcallback"

stripe:
	smee -u "https://smee.io/WeoKBRir10gZf0Lf" -t "http://localhost:7999/stripe/webhook"
