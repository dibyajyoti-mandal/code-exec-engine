#!/bin/sh
set -e

# Write Python code to file
echo "$CODE" > code.py

# Run with a 3-second time limit
timeout 3 python3 code.py
