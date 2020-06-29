#!/usr/bin/env bash
# Launches unit tests
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/posmining/keeper
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/posmining/types
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/structure/keeper
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/coins/keeper