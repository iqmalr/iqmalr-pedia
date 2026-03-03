#!/usr/bin/env bash

# Jalankan semua service Go dengan air secara paralel

echo "🚀 Starting go-product Go services with Air..."
# go-product service
(
  cd go-product/ || exit
  echo "▶️ Running go-product service..."
  air
) &

wait