#!/usr/bin/env bash
set -euo pipefail

make clean
make build-linux
make build-darwin
make build-windows
