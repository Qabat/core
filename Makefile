#!/usr/bin/env make
VER = v0.3.3
BUILD = $(shell git rev-parse --short HEAD)
FULL_VER = $(VER)-$(BUILD)

GOCMD=./cmd
ifeq ($(GO), )
    GO=go
endif

ifeq ($(GOPATH), )
    GOPATH=$(shell ls -d ~/go)
endif

INSTALLDIR=${GOPATH}/bin/


TARGETDIR=target
ARCH := $(shell uname -m)
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
OS_ARCH := $(OS)_$(ARCH)

WORKER=${TARGETDIR}/sonmworker_$(OS_ARCH)
CLI=${TARGETDIR}/sonmcli_$(OS_ARCH)
LOCAL_NODE=${TARGETDIR}/sonmnode_$(OS_ARCH)
AUTOCLI=${TARGETDIR}/autocli_$(OS_ARCH)
DWH=${TARGETDIR}/sonmdwh_$(OS_ARCH)
RENDEZVOUS=${TARGETDIR}/sonmrendezvous_$(OS_ARCH)
RELAY=${TARGETDIR}/sonmrelay_$(OS_ARCH)
LSGPU=${TARGETDIR}/lsgpu_$(OS_ARCH)
PANDORA=${TARGETDIR}/pandora_$(OS_ARCH)

TAGS=nocgo

GPU_SUPPORT?=false
ifeq ($(GPU_SUPPORT),true)
    GPU_TAGS=cl
    # required for build nvidia-docker libs with NVML included via cgo
    NV_CGO=vendor/github.com/sshaman1101/nvidia-docker/build
    CGO_LDFLAGS=-L$(shell pwd)/${NV_CGO}/lib
    CGO_CFLAGS=-I$(shell pwd)/${NV_CGO}/include
    CGO_LDFLAGS_ALLOW='-Wl,--unresolved-symbols=ignore-in-object-files'
endif


ifeq ($(OS),linux)
SED=sed -i 's/github\.com\/sonm-io\/core\/vendor\///g' insonmnia/dealer/hub_mock.go
endif

ifeq ($(OS),darwin)
SED=sed -i "" 's/github\.com\/sonm-io\/core\/vendor\///g' insonmnia/dealer/hub_mock.go
endif

LDFLAGS = -X main.appVersion=$(FULL_VER)

.PHONY: fmt vet test

all: mock vet fmt build test

build/worker:
	@echo "+ $@"
	CGO_LDFLAGS_ALLOW=${CGO_LDFLAGS_ALLOW} CGO_LDFLAGS=${CGO_LDFLAGS} CGO_CFLAGS=${CGO_CFLAGS} ${GO} build -tags "$(TAGS) $(GPU_TAGS)" -ldflags "-s $(LDFLAGS)" -o ${WORKER} ${GOCMD}/worker

build/dwh:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${DWH} ${GOCMD}/dwh

build/rv:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -o ${RENDEZVOUS} ${GOCMD}/rv

build/rendezvous: build/rv

build/relay:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -o ${RELAY} ${GOCMD}/relay

build/cli:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${CLI} ${GOCMD}/cli

build/node:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${LOCAL_NODE} ${GOCMD}/node

build/lsgpu:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${LSGPU} ${GOCMD}/lsgpu

build/cli_win32:
	@echo "+ $@"
	GOOS=windows GOARCH=386 ${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS).win32" -o ${CLI}_win32.exe ${GOCMD}/cli

build/node_win32:
	@echo "+ $@"
	GOOS=windows GOARCH=386 ${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS).win32" -o ${LOCAL_NODE}_win32.exe ${GOCMD}/node

build/autocli:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${AUTOCLI} ${GOCMD}/autocli

build/pandora:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${PANDORA} ${GOCMD}/pandora

build/pandora_linux:
	@echo "+ $@"
	GOOS=linux ${GO} build -tags "$(TAGS)" -ldflags "-s $(LDFLAGS)" -o ${TARGETDIR}/pandora_linux_x86_64 ${GOCMD}/pandora

build/insomnia: build/worker build/cli build/node

build/aux: build/relay build/rv build/dwh build/pandora

build: build/insomnia build/aux

install: all
	@echo "+ $@"
	mkdir -p ${INSTALLDIR}
	cp ${WORKER} ${CLI} ${LOCAL_NODE} ${INSTALLDIR}

vet:
	@echo "+ $@"
	@go tool vet $(shell ls -1 -d */ | grep -v -e vendor -e contracts)

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)

test: mock
	@echo "+ $@"
	${GO} test -tags nocgo $(shell go list ./... | grep -vE 'vendor|blockchain')

grpc:
	@echo "+ $@"
	@if ! which protoc > /dev/null; then echo "protoc protobuf compiler required for build"; exit 1; fi;
	@if ! which protoc-gen-grpccmd > /dev/null; then echo "protoc-gen-grpccmd protobuf plugin required for build.\nRun \`go get -u github.com/sshaman1101/grpccmd/cmd/protoc-gen-grpccmd\`"; exit 1; fi;
	@protoc -I proto proto/*.proto --grpccmd_out=proto/

build_mockgen:
	cd ./vendor/github.com/golang/mock/mockgen/ && go install

mock: build_mockgen
	mockgen -package miner -destination insonmnia/miner/overseer_mock.go -source insonmnia/miner/overseer.go
	mockgen -package task_config -destination cmd/cli/task_config/config_mock.go  -source cmd/cli/task_config/config.go
	mockgen -package accounts -destination accounts/keys_mock.go  -source accounts/keys.go
	mockgen -package benchmarks -destination insonmnia/benchmarks/benchmarks_mock.go  -source insonmnia/benchmarks/benchmarks.go
	mockgen -package blockchain -destination blockchain/api_mock.go  -source blockchain/api.go
	mockgen -package sonm -destination proto/marketplace_mock.go  -source proto/marketplace.pb.go
	mockgen -package sonm -destination proto/dwh_mock.go  -source proto/dwh.pb.go
	mockgen -package config -destination cmd/cli/config/config_mock.go  -source cmd/cli/config/config.go \
		-aux_files accounts=accounts/keys.go,logging=insonmnia/logging/logging.go

clean:
	rm -f ${WORKER} ${CLI} ${LOCAL_NODE} ${AUTOCLI} ${RENDEZVOUS}
	find . -name "*_mock.go" | xargs rm -f

deb:
	debuild --no-lintian --preserve-env -uc -us -i -I -b
	debuild clean

coverage:
	.ci/coverage.sh
