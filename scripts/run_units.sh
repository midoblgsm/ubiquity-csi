#!/usr/bin/env bash
echo "Setting up ginkgo and gomega"
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega



echo "Run unit tests"
ginkgo -r --skip vendor -ldflags -s
