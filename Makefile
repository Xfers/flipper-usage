default: build

flipper-usage-scan:
	go build -ldflags="-s -w" -o build/flipper-usage-scan cmd/flipper-usage-scan/main.go
flipper-merge-result:
	go build -ldflags="-s -w" -o build/flipper-merge-result cmd/flipper-merge-result/main.go

build: flipper-usage-scan flipper-merge-result

install: build
