#!/usr/bin/env bash

# Генерируем весь state с самого нуля
rm -rf ~/.ouroborosd/ # Выпиливаем все для старта с чистого листа
rm -rf ~/.ouroboroscli/ # Выпиливаем все для старта с чистого листа

# Инициализируем блокчейн
ouroborosd init genesis --chain-id ouroboros

# Генерим стандартные аккаунты в expect скриптах
#./scripts/generate_jack.exp
#./scripts/generate_alice.exp

# Добавляем первый аккаунт в генезис
ouroborosd add-genesis-account $(appcli keys show admin-wallet@ouroboros-crypto.com -a) 10000000000000ouro,10000000000000stake

# Конфигурируем cli
ouroboroscli config chain-id ouroboros
ouroboroscli config output json
ouroboroscli config indent true
ouroboroscli config trust-node true

# Генерируем стартовую транзакцию и пуляем ее в validators
echo "abcdef123" | ouroborosd gentx --name admin-wallet@ouroboros-crypto.com --amount 100000000stake

ouroborosd collect-gentxs