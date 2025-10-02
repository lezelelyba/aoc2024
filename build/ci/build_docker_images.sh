#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( realpath "$SCRIPT_DIR/../.." )"

cd "$PROJECT_ROOT"

docker build -f build/ci/docker/web.Dockerfile -t advent2024.web .
docker build -f build/ci/docker/cli.Dockerfile -t advent2024.cli .
