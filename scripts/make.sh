#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

cd "$(git rev-parse --show-toplevel)"

target="${1:-help}"

for f in cmd/*/Makefile; do
  (
    cd "$(dirname "${f}")"
    echo -e "\n$(pwd):"
    make "${target}"
  )
done
