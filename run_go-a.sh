#!/usr/bin/env bash

# Jalankan semua service Go dengan air secara paralel

echo "🚀 Starting auth Go services with Air..."

# go-auth service
(
  cd go-auth/v2 || exit
  echo "▶️ Running go-auth service..."
  air
) &

# Tunggu semua proses selesai
wait
