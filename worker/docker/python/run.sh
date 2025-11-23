#!/bin/sh
set -e

echo "$CODE" > code.js

timeout 3 node code.js
