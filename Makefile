# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: mit android ios mit-cross swarm evm all test clean
.PHONY: mit-linux mit-linux-386 mit-linux-amd64 mit-linux-mips64 mit-linux-mips64le
.PHONY: mit-linux-arm mit-linux-arm-5 mit-linux-arm-6 mit-linux-arm-7 mit-linux-arm64
.PHONY: mit-darwin mit-darwin-386 mit-darwin-amd64
.PHONY: mit-windows mit-windows-386 mit-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

mit:
	build/env.sh go run build/ci.go install ./cmd/mit
	@echo "Done building."
	@echo "Run \"$(GOBIN)/mit\" to launch mit."

createpx509:
	build/env.sh go run build/ci.go install ./cmd/createpx509
	@echo "Done building."
	@echo "Run \"$(GOBIN)/createpx509\" to use createpx509."

readx509:
	build/env.sh go run build/ci.go install ./cmd/readx509
	@echo "Done building."
	@echo "Run \"$(GOBIN)/readx509\" to use readx509."

x509cert:
	build/env.sh go run build/ci.go install ./cmd/x509cert
	@echo "Done building."
	@echo "Run \"$(GOBIN)/x509cert\" to use x509cert."

puppeth:
	build/env.sh go run build/ci.go install ./cmd/puppeth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/puppeth\" to launch puppeth."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/mit.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Mit.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

mit-cross: mit-linux mit-darwin mit-windows mit-android mit-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/mit-*

mit-linux: mit-linux-386 mit-linux-amd64 mit-linux-arm mit-linux-mips64 mit-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-*

mit-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/mit
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep 386

mit-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/mit
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep amd64

mit-linux-arm: mit-linux-arm-5 mit-linux-arm-6 mit-linux-arm-7 mit-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep arm

mit-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/mit
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep arm-5

mit-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/mit
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep arm-6

mit-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/mit
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep arm-7

mit-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/mit
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep arm64

mit-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/mit
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep mips

mit-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/mit
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep mipsle

mit-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/mit
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep mips64

mit-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/mit
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/mit-linux-* | grep mips64le

mit-darwin: mit-darwin-386 mit-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/mit-darwin-*

mit-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/mit
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/mit-darwin-* | grep 386

mit-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/mit
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/mit-darwin-* | grep amd64

mit-windows: mit-windows-386 mit-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/mit-windows-*

mit-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/mit
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/mit-windows-* | grep 386

mit-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/mit
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/mit-windows-* | grep amd64
