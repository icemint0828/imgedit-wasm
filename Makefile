.PHONY: all build check format test vet
all: build
check: format vet test

build: check
	
	@echo "[+] Building the wasm version"
	@GOOS=js GOARCH=wasm go build -o docs/main.wasm wasm/main.go

	@echo "[+] Done"

format:
	@echo "[+] Formatting files"
	@gofmt -w *.go

vet:
	@echo "[+] Running Go vet"
	@go vet

test:
	@echo "[+] Running tests"
	@go test