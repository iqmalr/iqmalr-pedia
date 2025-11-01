#!/usr/bin/env bash

# Jalankan semua service Go dengan air secara paralel

echo "🚀 Starting go-api-gateway Go services with Air..."
# go-api-gateway service
(
  cd go-api-gateway/v2 || exit
  echo "▶️ Running go-api-gateway service..."
  air
) &

wait