#!/usr/bin/env bash
set -euo pipefail

make fmt
make vet
make test
make release
