#!/usr/bin/env bash

# Jalankan semua service Go dengan air secara paralel

echo "🚀 Starting go-vendors Go services with Air..."
# go-vendors service
(
  cd go-vendors/ || exit
  echo "▶️ Running go-vendors service..."
  air
) &

wait