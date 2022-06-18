.PHONY: all build check format test vet
all: build
check: format vet test

build: check
	
	@echo "[+] Building the wasm version"
	@GOOS=js GOARCH=wasm go build -o docs/main.wasm main.go

	@echo "[+] Done"

format:
	@echo "[+] Formatting files"
	@GOOS=js GOARCH=wasm gofmt -w *.go

vet:
	@echo "[+] Running Go vet"
	@GOOS=js GOARCH=wasm go vet

test:
	@echo "[+] Running tests"
	@GOOS=js GOARCH=wasm go test