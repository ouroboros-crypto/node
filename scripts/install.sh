#!/usr/bin/env bash
# Устанавливаем все необходимое
GO111MODULE=on go get
GO111MODULE=on go install ./cmd/ouroborosd
GO111MODULE=on go install ./cmd/ouroboroscli

if [ ! -f "~/.ouroborosd/config/node_key.json" ]
then
  mkdir -p ~/.ouroborosd/config
  mkdir -p ~/.ouroborosd/data

  cp -r ./installation/genesis.json ~/.ouroborosd/config/
  cp -r ./installation/config.toml ~/.ouroborosd/config/

  ouroboroscli config chain-id ouroboros
  ouroboroscli config output json
  ouroboroscli config indent true
  ouroboroscli config trust-node true
  ouroboroscli config node tcp://127.0.0.1:26657
fi

if [ ! -f "~/.ouroborosd/config/genesis.json" ]
then
  cp -r ./installation/genesis.json ~/.ouroborosd/config/
fi

if [ ! -f "~/.ouroborosd/config/config.toml" ]
then
  cp -r ./installation/config.toml ~/.ouroborosd/config/
fi
