#!/bin/bash
set -e

cd "$(dirname "$0")/.."

if [ ! -d "venv" ]; then
    echo "venv not found, run: python3 -m venv venv && venv/bin/pip install -r requirements-dev.txt"
    exit 1
fi

venv/bin/pytest tests/ -v
