BIN     := gsm
CMD     := ./cmd/gsm
LDFLAGS := -ldflags="-s -w"

.PHONY: build test fmt vet clean install

build:
	go build $(LDFLAGS) -o $(BIN) $(CMD)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BIN)

install:
	go install $(LDFLAGS) $(CMD)
