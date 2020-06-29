#!/usr/bin/env bash
# Installs all the required things to your system
GO111MODULE=on go get
GO111MODULE=on go install ./cmd/ouroboroscli
GO111MODULE=on go install ./cmd/ouroborosd
